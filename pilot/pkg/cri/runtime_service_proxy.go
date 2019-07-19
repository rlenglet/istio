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
	"time"

	"go.uber.org/zap"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	// PollingPeriod is the period for polling pod sandbox and container statuses.
	PollingPeriod = 500 * time.Millisecond
)

// ExtendedRuntimeServiceServer extends RuntimeServiceServer with higher-level utility functions.
type ExtendedRuntimeServiceServer interface {
	api.RuntimeServiceServer

	// WaitForPodSandboxStatus polls the status of the given pod sandbox until it reaches a desired
	// status or the context is cancelled, whichever happens first.
	WaitForPodSandboxStatus(ctx context.Context, logFields []zap.Field,
		podSandboxID string, statusPredicate func(status *api.PodSandboxStatus) (bool, error)) (
		*api.PodSandboxStatus, error)

	// CreateAndStartContainer creates and starts a container in the given pod sandbox and with the
	// given configuration, and returns the container ID if it has been created.
	CreateAndStartContainer(ctx context.Context, logFields []zap.Field,
		podSandboxID string, podSandboxConfig *api.PodSandboxConfig,
		containerConfig *api.ContainerConfig) (string, error)

	// WaitForContainerStatus polls the status of the given container until it reaches a desired
	// status or the context is cancelled, whichever happens first.
	WaitForContainerStatus(ctx context.Context, logFields []zap.Field,
		containerID string, statusPredicate func(status *api.ContainerStatus) (bool, error)) (
		*api.ContainerStatus, error)
}

// RuntimeServiceProxy is a RuntimeServiceServer that forwards all received calls to a backend
// RuntimeServiceClient. This is intended to be overridden to intercept specific calls.
type RuntimeServiceProxy struct {
	// runtimeService is the backend RuntimeServiceClient that calls are forwarded to.
	runtimeService api.RuntimeServiceClient

	// runtimeServiceLogger logs calls forwarded to runtimeService.
	runtimeServiceLogger Logger
}

var _ ExtendedRuntimeServiceServer = &RuntimeServiceProxy{}

// NewRuntimeServiceProxy creates a new RuntimeServiceProxy backed by the given RuntimeServiceClient.
func NewRuntimeServiceProxy(runtimeService api.RuntimeServiceClient, runtimeServiceLogger Logger) *RuntimeServiceProxy {
	return &RuntimeServiceProxy{
		runtimeService:       runtimeService,
		runtimeServiceLogger: runtimeServiceLogger,
	}
}

func (p *RuntimeServiceProxy) WaitForPodSandboxStatus(ctx context.Context, logFields []zap.Field,
	podSandboxID string, statusPredicate func(status *api.PodSandboxStatus) (bool, error)) (
	*api.PodSandboxStatus, error) {
	log.Debug("Waiting for pod sandbox to reach desired status", logFields...)

	statusReq := &api.PodSandboxStatusRequest{PodSandboxId: podSandboxID}
	for {
		statusRes, err := p.PodSandboxStatus(ctx, statusReq)
		if err != nil {
			return nil, fmt.Errorf("failed to get pod sandbox status: %v", err)
		}

		status := statusRes.GetStatus()

		done, err := statusPredicate(status)
		if err != nil {
			return nil, err
		}
		if done {
			return status, nil
		}

		// Sleep before checking again.
		t := time.NewTimer(PollingPeriod)
		select {
		case <-ctx.Done():
			t.Stop()
			return nil, ctx.Err()
		case <-t.C:
		}
	}
}

// PodSandboxIsReadyAndHasIP is a predicate to use with WaitForPodSandboxStatus which returns
// whether the pod is in ready state and has been assigned an IP address.
func PodSandboxIsReadyAndHasIP(status *api.PodSandboxStatus) (bool, error) {
	return status.GetState() == api.PodSandboxState_SANDBOX_READY &&
		status.GetNetwork().GetIp() != "", nil
}

func (p *RuntimeServiceProxy) CreateAndStartContainer(ctx context.Context, logFields []zap.Field,
	podSandboxID string, podSandboxConfig *api.PodSandboxConfig,
	containerConfig *api.ContainerConfig) (string, error) {
	log.Debug("Creating container", logFields...)

	createReq := &api.CreateContainerRequest{
		PodSandboxId:  podSandboxID,
		Config:        containerConfig,
		SandboxConfig: podSandboxConfig,
	}
	createRes, err := p.CreateContainer(ctx, createReq)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}
	containerID := createRes.GetContainerId()

	log.Debug("Starting container", logFields...)

	startReq := &api.StartContainerRequest{ContainerId: containerID}
	_, err = p.StartContainer(ctx, startReq)
	if err != nil {
		return containerID, fmt.Errorf("failed to start container: %v", err)
	}

	return containerID, nil
}

func (p *RuntimeServiceProxy) WaitForContainerStatus(ctx context.Context, logFields []zap.Field,
	containerID string, statusPredicate func(status *api.ContainerStatus) (bool, error)) (
	*api.ContainerStatus, error) {
	log.Debug("Waiting for container to reach desired status", logFields...)

	statusReq := &api.ContainerStatusRequest{ContainerId: containerID}
	for {
		statusRes, err := p.ContainerStatus(ctx, statusReq)
		if err != nil {
			return nil, fmt.Errorf("failed to get container status: %v", err)
		}

		status := statusRes.GetStatus()

		done, err := statusPredicate(status)
		if err != nil {
			return nil, err
		}
		if done {
			return status, nil
		}

		// Sleep before checking again.
		t := time.NewTimer(PollingPeriod)
		select {
		case <-ctx.Done():
			t.Stop()
			return nil, ctx.Err()
		case <-t.C:
		}
	}
}

func (p *RuntimeServiceProxy) Version(ctx context.Context, req *api.VersionRequest) (*api.VersionResponse, error) {
	res, err := p.runtimeService.Version(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) RunPodSandbox(ctx context.Context, req *api.RunPodSandboxRequest) (*api.RunPodSandboxResponse, error) {
	res, err := p.runtimeService.RunPodSandbox(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) StopPodSandbox(ctx context.Context, req *api.StopPodSandboxRequest) (*api.StopPodSandboxResponse, error) {
	res, err := p.runtimeService.StopPodSandbox(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) RemovePodSandbox(ctx context.Context, req *api.RemovePodSandboxRequest) (*api.RemovePodSandboxResponse, error) {
	res, err := p.runtimeService.RemovePodSandbox(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) PodSandboxStatus(ctx context.Context, req *api.PodSandboxStatusRequest) (*api.PodSandboxStatusResponse, error) {
	res, err := p.runtimeService.PodSandboxStatus(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ListPodSandbox(ctx context.Context, req *api.ListPodSandboxRequest) (*api.ListPodSandboxResponse, error) {
	res, err := p.runtimeService.ListPodSandbox(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) CreateContainer(ctx context.Context, req *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	res, err := p.runtimeService.CreateContainer(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) StartContainer(ctx context.Context, req *api.StartContainerRequest) (*api.StartContainerResponse, error) {
	res, err := p.runtimeService.StartContainer(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) StopContainer(ctx context.Context, req *api.StopContainerRequest) (*api.StopContainerResponse, error) {
	res, err := p.runtimeService.StopContainer(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) RemoveContainer(ctx context.Context, req *api.RemoveContainerRequest) (*api.RemoveContainerResponse, error) {
	res, err := p.runtimeService.RemoveContainer(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ListContainers(ctx context.Context, req *api.ListContainersRequest) (*api.ListContainersResponse, error) {
	res, err := p.runtimeService.ListContainers(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ContainerStatus(ctx context.Context, req *api.ContainerStatusRequest) (*api.ContainerStatusResponse, error) {
	res, err := p.runtimeService.ContainerStatus(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) UpdateContainerResources(ctx context.Context, req *api.UpdateContainerResourcesRequest) (*api.UpdateContainerResourcesResponse, error) {
	res, err := p.runtimeService.UpdateContainerResources(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ReopenContainerLog(ctx context.Context, req *api.ReopenContainerLogRequest) (*api.ReopenContainerLogResponse, error) {
	res, err := p.runtimeService.ReopenContainerLog(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ExecSync(ctx context.Context, req *api.ExecSyncRequest) (*api.ExecSyncResponse, error) {
	res, err := p.runtimeService.ExecSync(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) Exec(ctx context.Context, req *api.ExecRequest) (*api.ExecResponse, error) {
	res, err := p.runtimeService.Exec(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) Attach(ctx context.Context, req *api.AttachRequest) (*api.AttachResponse, error) {
	res, err := p.runtimeService.Attach(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) PortForward(ctx context.Context, req *api.PortForwardRequest) (*api.PortForwardResponse, error) {
	res, err := p.runtimeService.PortForward(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ContainerStats(ctx context.Context, req *api.ContainerStatsRequest) (*api.ContainerStatsResponse, error) {
	res, err := p.runtimeService.ContainerStats(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) ListContainerStats(ctx context.Context, req *api.ListContainerStatsRequest) (*api.ListContainerStatsResponse, error) {
	res, err := p.runtimeService.ListContainerStats(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) UpdateRuntimeConfig(ctx context.Context, req *api.UpdateRuntimeConfigRequest) (*api.UpdateRuntimeConfigResponse, error) {
	res, err := p.runtimeService.UpdateRuntimeConfig(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *RuntimeServiceProxy) Status(ctx context.Context, req *api.StatusRequest) (*api.StatusResponse, error) {
	res, err := p.runtimeService.Status(ctx, req)
	p.runtimeServiceLogger.LogCall(req, res, err)
	return res, err
}
