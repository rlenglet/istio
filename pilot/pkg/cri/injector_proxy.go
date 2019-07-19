// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cri

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

//const (
//	// PodAnnotationKeyInject is the key of a pod annotation that can be set to disable Istio proxy
//	// injection.
//	PodAnnotationKeyInject = "sidecar.istio.io/inject"
//
//	// PodAnnotationKeyStatus is the key of the pod annotation added by Istio's pod injector.
//	PodAnnotationKeyStatus = "sidecar.istio.io/status"
//
//	// PodAnnotationStatusKeyContainers is the field name in the JSON object stored in the
//	// sidecar.istio.io/status annotation which contains the list of injected containers.
//	PodAnnotationStatusKeyContainers = "containers"
//)

// InjectorProxy is a RuntimeServiceProxy and ImageServiceProxy that injects Istio proxies into
// pods created via CRI.
type InjectorProxy struct {
	RuntimeServiceProxy
	ImageServiceProxy

	// injectionConfigurator generates Istio pod injection configurations.
	injectionConfigurator InjectionConfigurator

	// excludedNamespaces is the set of k8s namespaces excluded from Istio proxy injection.
	excludedNamespaces map[string]struct{}
}

var _ api.RuntimeServiceServer = &InjectorProxy{}
var _ api.ImageServiceServer = &InjectorProxy{}

// NewInjectorProxy creates a new InjectorProxy backed by the given RuntimeServiceClient and
// ImageServiceClient.
func NewInjectorProxy(runtimeService api.RuntimeServiceClient,
	imageService api.ImageServiceClient, logger Logger,
	injectionConfigurator InjectionConfigurator, excludedNamespaces []string) *InjectorProxy {
	var excludedNamespacesMap map[string]struct{}
	if len(excludedNamespaces) > 1 {
		excludedNamespacesMap = make(map[string]struct{}, len(excludedNamespaces))
		for _, ns := range excludedNamespaces {
			excludedNamespacesMap[ns] = struct{}{}
		}
	}

	return &InjectorProxy{
		RuntimeServiceProxy: RuntimeServiceProxy{
			runtimeService:       runtimeService,
			runtimeServiceLogger: logger,
		},
		ImageServiceProxy: ImageServiceProxy{
			imageService:       imageService,
			imageServiceLogger: logger,
		},
		injectionConfigurator: injectionConfigurator,
		excludedNamespaces:    excludedNamespacesMap,
	}
}

func (p *InjectorProxy) RunPodSandbox(ctx context.Context, req *api.RunPodSandboxRequest) (
	*api.RunPodSandboxResponse, error) {
	logFields := []zap.Field{
		zap.String("pod_name", req.GetConfig().GetMetadata().GetName()),
		zap.String("pod_namespace", req.GetConfig().GetMetadata().GetNamespace()),
	}

	log.Info("Generating Istio proxy injection configuration", logFields...)

	podInjectionConfigurator, err := p.injectionConfigurator.ConfigurePodInjection(ctx, logFields,
		req.GetConfig())
	if err != nil {
		log.Error("Error while generating Istio proxy injection configuration",
			WithError(err, logFields...)...)
		return nil, fmt.Errorf("failed while deciding whether to inject Istio proxy into pod: %v",
			err)
	}
	injectionEnabled := podInjectionConfigurator != nil

	// If the istio-proxy container will be injected, update the pod config with DNS config and
	// ports.
	if injectionEnabled {
		// TODO(rlenglet): append the DNSConfig.

		req.GetConfig().PortMappings = AppendPortMappings(req.GetConfig().PortMappings,
			podInjectionConfigurator.PortMappings()...)
	}

	log.Info("Running pod", logFields...)

	res, err := p.RuntimeServiceProxy.RunPodSandbox(ctx, req)
	if err != nil {
		return res, err
	}

	podSandboxID := res.PodSandboxId
	podSandboxIDField := zap.String("pod_id", podSandboxID)

	logFields = []zap.Field{
		podSandboxIDField,
		zap.String("pod_name", req.GetConfig().GetMetadata().GetName()),
		zap.String("pod_namespace", req.GetConfig().GetMetadata().GetNamespace()),
	}

	if !injectionEnabled {
		log.Info("Skipping pod Istio proxy injection", logFields...)
		return res, nil
	}

	log.Info("Waiting for pod to be ready", logFields...)

	// Wait until the pod sandbox:
	// - is in the ready state,
	// - has an IP address allocated.
	podSandboxStatus, err := p.WaitForPodSandboxStatus(ctx, logFields,
		podSandboxID, PodSandboxIsReadyAndHasIP)
	if err != nil {
		log.Error("Error while waiting for pod to be ready", WithError(err, logFields...)...)
		// TODO(rlenglet): Remove the sandbox in case of failure.
		return nil, err
	}

	log.Debug("Pod was allocated IP address",
		WithField(zap.String("pod_ip", podSandboxStatus.GetNetwork().GetIp()), logFields...)...)

	log.Info("Configuring traffic redirection to Istio proxy", logFields...)

	podInjectionConfig, err := podInjectionConfigurator.ConfigurePodInjectionForStatus(logFields,
		podSandboxStatus)
	if err != nil {
		log.Error("Error while generating Istio proxy injection configuration",
			WithError(err, logFields...)...)
		return nil, fmt.Errorf("failed while configuring Istio proxy injection into pod: %v", err)
	}

	if err := p.RunIstioInitContainer(ctx, logFields, podSandboxID, req.Config,
		&podInjectionConfig.IstioInitContainer); err != nil {
		log.Error("Error while configuring traffic redirection to Istio proxy",
			WithError(err, logFields...)...)
		// TODO(rlenglet): Remove the sandbox in case of failure.
		return nil, err
	}

	if !podInjectionConfig.IstioProxyInjectionEnabled {
		return res, nil
	}

	log.Info("Starting Istio proxy", logFields...)

	if err := p.StartIstioProxyContainer(ctx, logFields, podSandboxID, req.Config,
		&podInjectionConfig.IstioProxyContainer); err != nil {
		log.Error("Error while starting Istio proxy", WithError(err, logFields...)...)
		// TODO(rlenglet): Stop and remove the sandbox in case of failure.
		return nil, err
	}

	return res, nil
}

//// InjectionEnabled returns whether the pod sandbox with the given configuration has been or is to
//// be injected with the Istio proxy.
//func (p *InjectorProxy) InjectionEnabled(podSandboxConfig *api.PodSandboxConfig) bool {
//	logFields := []zap.Field{
//		zap.String("pod_name", podSandboxConfig.GetMetadata().GetName()),
//		zap.String("pod_namespace", podSandboxConfig.GetMetadata().GetNamespace()),
//	}
//
//	enabled := true
//
//	if _, found := p.excludedNamespaces[podSandboxConfig.GetMetadata().GetNamespace()]; found {
//		log.Info("Pod Istio proxy injection disabled for pod via namespace exclusion", logFields...)
//		enabled = false
//	}
//
//	if val, found := podSandboxConfig.GetAnnotations()[PodAnnotationKeyInject]; found {
//		if injectEnabled, err := strconv.ParseBool(val); !injectEnabled && err == nil {
//			log.Info("Pod Istio proxy injection disabled via "+PodAnnotationKeyInject+
//				" pod annotation", logFields...)
//			enabled = false
//		}
//	}
//
//	return enabled
//}
//
//// ProxyContainerEnabled returns whether the pod sandbox with the given configuration is to be
//// injected with the Istio proxy.
//func (p *InjectorProxy) ProxyContainerEnabled(podSandboxConfig *api.PodSandboxConfig) bool {
//	logFields := []zap.Field{
//		zap.String("pod_name", podSandboxConfig.GetMetadata().GetName()),
//		zap.String("pod_namespace", podSandboxConfig.GetMetadata().GetNamespace()),
//	}
//
//	// Check in the injector status whether an istio-proxy container has already been injected.
//	if status, found := podSandboxConfig.GetAnnotations()[PodAnnotationKeyStatus]; found {
//		var injectionStatus map[string]interface{}
//		err := json.Unmarshal([]byte(status), &injectionStatus)
//		if err != nil {
//			// Annotation can't be parsed as a JSON object. Disable injection to be safe.
//			log.Error("Pod Istio proxy injection disabled due to error parsing "+
//				PodAnnotationKeyStatus+" pod annotation", logFields...)
//			return false
//		}
//
//		if injectionStatus[PodAnnotationStatusKeyContainers] == nil {
//			// Not already injected.
//			return true
//		}
//
//		containers, ok := injectionStatus[PodAnnotationStatusKeyContainers].([]interface{})
//		if !ok {
//			// "containers" can't be parsed as a JSON list. Disable injection to be safe.
//			log.Error("Pod Istio proxy injection disabled due to error parsing "+
//				PodAnnotationKeyStatus+" pod annotation's "+PodAnnotationStatusKeyContainers+
//				" field", logFields...)
//			return false
//		}
//
//		for _, container := range containers {
//			if container == "istio-proxy" {
//				log.Info("Pod Istio proxy injection disabled since it is already injected",
//					logFields...)
//				return false
//			}
//		}
//	}
//
//	return true
//}

// CreateContainerLogPath creates the logging parent directory for a container and returns its log
// path.
func CreateContainerLogPath(logFields []zap.Field, podSandboxConfig *api.PodSandboxConfig,
	containerMetadata *api.ContainerMetadata) (string, error) {
	// kubelet's generateContainerConfig() determines the log directory as:
	//     filepath.Join(podSandboxConfig.GetLogDirectory(), containerMetadata.GetName())
	//
	// Add __istio__ into the path to reduce the risk of collisions.
	containerLogDir := filepath.Join("__istio__", containerMetadata.GetName())
	logDir := filepath.Join(podSandboxConfig.GetLogDirectory(), containerLogDir)

	log.Debug("Creating container log directory", logFields...)

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create container log directory %q: %v", logDir, err)
	}

	return filepath.Join(containerLogDir,
		fmt.Sprintf("%d.log", containerMetadata.GetAttempt())), nil
}

func (p *InjectorProxy) StartInjectedContainer(ctx context.Context, logFields []zap.Field,
	podSandboxID string, podSandboxConfig *api.PodSandboxConfig,
	injectConfig *InjectedContainerConfig) (string, error) {
	// TODO(rlenglet): Set the auth config.
	imageRef, err := p.ImageRef(ctx, WithField(zap.String("image", injectConfig.Image), logFields...),
		podSandboxConfig, injectConfig.Image, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get reference of %q container image %q: %v",
			injectConfig.Name, injectConfig.Image, err)
	}

	containerMetadata := &api.ContainerMetadata{
		Name:    injectConfig.Name,
		Attempt: 0,
	}

	logPath, err := CreateContainerLogPath(logFields, podSandboxConfig, containerMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to create %q container log directory: %v",
			injectConfig.Name, err)
	}

	// Add an annotation to identify the containers managed by this proxy.
	annotations := AddIstioManagedAnnotation(injectConfig.Annotations)

	containerConfig := &api.ContainerConfig{
		Metadata: containerMetadata,
		Image: &api.ImageSpec{
			Image: imageRef,
		},
		Command:     injectConfig.Command,
		Args:        injectConfig.Args,
		WorkingDir:  injectConfig.WorkingDir,
		Envs:        injectConfig.Envs,
		Mounts:      injectConfig.Mounts,
		Labels:      injectConfig.Labels,
		Annotations: annotations,
		LogPath:     logPath,
		Linux:       injectConfig.Linux,
	}

	containerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	containerID, err := p.CreateAndStartContainer(containerCtx, logFields, podSandboxID,
		podSandboxConfig, containerConfig)
	cancel()
	return containerID, err
}

// RunIstioInitContainer runs the istio-init container in the given sandbox and waits until it
// exits.
func (p *InjectorProxy) RunIstioInitContainer(ctx context.Context, logFields []zap.Field,
	podSandboxID string, podSandboxConfig *api.PodSandboxConfig,
	injectConfig *InjectedContainerConfig) error {
	logFields = WithField(zap.String("container_name", injectConfig.Name), logFields...)
	
	containerID, err := p.StartInjectedContainer(ctx, logFields, podSandboxID, podSandboxConfig, injectConfig)
	defer func() {
		log.Debug("Removing container", logFields...)

		removeReq := &api.RemoveContainerRequest{ContainerId: containerID}
		_, _ = p.RuntimeServiceProxy.RemoveContainer(ctx, removeReq)
	}()
	if err != nil {
		return err
	}

	_, err = p.WaitForContainerStatus(ctx, logFields, containerID,
		func(status *api.ContainerStatus) (bool, error) {
			return status.GetState() == api.ContainerState_CONTAINER_EXITED, nil
		})
	if err != nil {
		return fmt.Errorf("timed out while waiting for %q container to exit: %v",
			injectConfig.Name, err)
	}

	// TODO(rlenglet): Log status, etc.
	// TODO(rlenglet): Collect logs from the container.

	return nil
}

// StartIstioProxyContainer runs the istio-proxy container in the given sandbox and waits until it
// is ready.
func (p *InjectorProxy) StartIstioProxyContainer(ctx context.Context, logFields []zap.Field,
	podSandboxID string, podSandboxConfig *api.PodSandboxConfig,
	injectConfig *InjectedContainerConfig) error {
	logFields = WithField(zap.String("container_name", injectConfig.Name), logFields...)

	containerID, err := p.StartInjectedContainer(ctx, logFields, podSandboxID, podSandboxConfig, injectConfig)
	if err != nil {
		return err
	}

	status, err := p.WaitForContainerStatus(ctx, logFields, containerID,
		func(status *api.ContainerStatus) (bool, error) {
			switch status.GetState() {
			case api.ContainerState_CONTAINER_RUNNING, api.ContainerState_CONTAINER_EXITED:
				return true, nil
			default:
				return false, nil
			}
		})
	if err != nil {
		return fmt.Errorf("timed out while waiting for %q container to run: %v",
			injectConfig.Name, err)
	}
	if status.GetState() != api.ContainerState_CONTAINER_RUNNING {
		return fmt.Errorf("%q container is in unexpected state %s",
			injectConfig.Name, status.GetState())
	}

	// TODO(rlenglet): Implement readiness probe.

	// TODO(rlenglet): Log status, etc.
	// TODO(rlenglet): Collect logs from the container.

	return nil
}

func (p *InjectorProxy) CreateContainer(ctx context.Context, req *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	// Error out if the container contains an annotation that indicates that it's managed by this
	// proxy.
	if ContainsIstioManagedAnnotation(req.GetConfig().GetAnnotations()) {
		log.Error("Rejecting attempt to create container with reserved annotation",
			zap.String("pod_id", req.GetPodSandboxId()),
			zap.String("pod_name", req.GetSandboxConfig().GetMetadata().GetName()),
			zap.String("pod_namespace", req.GetSandboxConfig().GetMetadata().GetName()),
			zap.String("container_name", req.GetConfig().GetMetadata().GetName()),
			zap.Uint32("container_attempt", req.GetConfig().GetMetadata().GetAttempt()),
			zap.String("image", req.GetConfig().GetImage().GetImage()))
		return nil, fmt.Errorf("attempt to create container with reserved annotation")
	}

	return p.RuntimeServiceProxy.CreateContainer(ctx, req)
}

func (p *InjectorProxy) ListContainers(ctx context.Context, req *api.ListContainersRequest) (*api.ListContainersResponse, error) {
	res, err := p.RuntimeServiceProxy.ListContainers(ctx, req)
	if err != nil {
		return nil, err
	}

	// Hide any container that is managed by this CRI proxy, to prevent kubelet from deleting it.
	if len(res.Containers) > 0 {
		newContainers := make([]*api.Container, 0, len(res.Containers))
		for _, container := range res.Containers {
			if !ContainsIstioManagedAnnotation(container.GetAnnotations()) {
				newContainers = append(newContainers, container)
			}
		}
		res.Containers = newContainers
	}

	return res, nil
}

func (p *InjectorProxy) ListContainerStats(ctx context.Context, req *api.ListContainerStatsRequest) (*api.ListContainerStatsResponse, error) {
	res, err := p.RuntimeServiceProxy.ListContainerStats(ctx, req)
	if err != nil {
		return nil, err
	}

	// Hide any container that is managed by this CRI proxy, to prevent kubelet from deleting it.
	if len(res.Stats) > 0 {
		newStats := make([]*api.ContainerStats, 0, len(res.Stats))
		for _, containerStat := range res.Stats {
			if !ContainsIstioManagedAnnotation(containerStat.GetAttributes().GetAnnotations()) {
				newStats = append(newStats, containerStat)
			}
		}
		res.Stats = newStats
	}

	return res, nil
}
