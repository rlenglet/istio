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

	"go.uber.org/zap"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// InjectionConfigurator configures Istio injection for pods.
type InjectionConfigurator interface {
	// ConfigurePodInjection generates the configuration of Istio injection for one specific pod.
	// If it returns nil, Istio injection is disabled for the pod.
	//
	// Configuration is done in two steps.
	// 1) ConfigurePodInjection checks whether the pod should be injected, before the pod is
	//    created, and generates an injection configuration template for the pod based on the pod
	//    sandbox configuration.
	// 2) Calling ConfigurePodInjectionForStatus on the returned object, after the pod is ready,
	//    generates the complete injection configuration based on both the pod's configuration and
	//    its status.
	ConfigurePodInjection(ctx context.Context, logFields []zap.Field,
		podSandboxConfig *api.PodSandboxConfig) (PodInjectionConfigurator, error)
}

// PodInjectionConfigurator configures Istio injection for one specific pod.
type PodInjectionConfigurator interface {
	// DNSConfig returns the additional DNS config for the pod.
	DNSConfig() *api.DNSConfig

	// PortMappings returns the list of port mappings to add into the pod.
	PortMappings() []*api.PortMapping

	// ConfigurePodInjectionForStatus generates the configuration of Istio injection for the pod
	// based on its status.
	ConfigurePodInjectionForStatus(logFields []zap.Field, podSandboxStatus *api.PodSandboxStatus) (
		*PodInjectionConfig, error)
}

// PodInjectionConfig configures Istio injections into a pod.
type PodInjectionConfig struct {
	// IstioInitContainer is the configuration of the istio-init container to inject.
	IstioInitContainer InjectedContainerConfig `json:"istio_init_container"`

	// IstioProxyInjectionEnabled indicates whether the istio-proxy container is to be injected.
	// If false, it has already been injected into the pod, e.g. via manual injection or via the
	// webhook.
	IstioProxyInjectionEnabled bool `json:"istio_proxy_injection_enabled"`

	// IstioProxyContainer is the configuration of the istio-proxy container to inject.
	// If IstioProxyInjectionEnabled is false, this contains the configuration of the
	// already-injected istio-proxy container.
	IstioProxyContainer InjectedContainerConfig `json:"istio_proxy_container"`
}

// InjectedContainerConfig is a subset of ContainerConfig for configuring an injected container.
type InjectedContainerConfig struct {
	// Name is the name of the container.
	Name string `json:"name"`

	// Image is the name of the image of the container.
	Image string `json:"image"`

	// Command is the command to execute in the container (i.e., entrypoint for docker).
	Command []string `json:"command,omitempty"`

	// Args are the arguments for the Command (i.e., command for docker).
	Args []string `json:"args"`

	// WorkingDir is the current working directory of the command.
	WorkingDir string `json:"working_dir,omitempty"`

	// Envs is the list of environment variable to set in the container.
	Envs []*api.KeyValue `json:"envs,omitempty"`

	// Mounts are mounts for the container.
	Mounts []*api.Mount `json:"mounts,omitempty"`

	// Labels are key-value pairs that may be used to scope and select individual resources.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are unstructured key-value map that may be used by the kubelet to store and
	// retrieve arbitrary metadata.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Linux is the configuration specific to Linux containers.
	Linux *api.LinuxContainerConfig `json:"linux,omitempty"`
}
