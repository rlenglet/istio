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

	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RuntimeServiceServerMock is a super type of mock implementations of RuntimeServiceServer for
// testing purposes.
// This mock's methods should be overridden to return proper responses instead of returning errors.
type RuntimeServiceServerMock struct{}

var _ api.RuntimeServiceServer = &RuntimeServiceServerMock{}

func (m *RuntimeServiceServerMock) Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) RunPodSandbox(context.Context, *api.RunPodSandboxRequest) (*api.RunPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) StopPodSandbox(context.Context, *api.StopPodSandboxRequest) (*api.StopPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) RemovePodSandbox(context.Context, *api.RemovePodSandboxRequest) (*api.RemovePodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) PodSandboxStatus(context.Context, *api.PodSandboxStatusRequest) (*api.PodSandboxStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ListPodSandbox(context.Context, *api.ListPodSandboxRequest) (*api.ListPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) CreateContainer(context.Context, *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) StartContainer(context.Context, *api.StartContainerRequest) (*api.StartContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) StopContainer(context.Context, *api.StopContainerRequest) (*api.StopContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) RemoveContainer(context.Context, *api.RemoveContainerRequest) (*api.RemoveContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ListContainers(context.Context, *api.ListContainersRequest) (*api.ListContainersResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ContainerStatus(context.Context, *api.ContainerStatusRequest) (*api.ContainerStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) UpdateContainerResources(context.Context, *api.UpdateContainerResourcesRequest) (*api.UpdateContainerResourcesResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ReopenContainerLog(context.Context, *api.ReopenContainerLogRequest) (*api.ReopenContainerLogResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ExecSync(context.Context, *api.ExecSyncRequest) (*api.ExecSyncResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) Exec(context.Context, *api.ExecRequest) (*api.ExecResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) Attach(context.Context, *api.AttachRequest) (*api.AttachResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) PortForward(context.Context, *api.PortForwardRequest) (*api.PortForwardResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ContainerStats(context.Context, *api.ContainerStatsRequest) (*api.ContainerStatsResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) ListContainerStats(context.Context, *api.ListContainerStatsRequest) (*api.ListContainerStatsResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) UpdateRuntimeConfig(context.Context, *api.UpdateRuntimeConfigRequest) (*api.UpdateRuntimeConfigResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceServerMock) Status(context.Context, *api.StatusRequest) (*api.StatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

// ImageServiceServerMock is a super type of mock implementations of ImageServiceServer for
// testing purposes.
// This mock's methods should be overridden to return proper responses instead of returning errors.
type ImageServiceServerMock struct{}

var _ api.ImageServiceServer = &ImageServiceServerMock{}

func (m *ImageServiceServerMock) ListImages(context.Context, *api.ListImagesRequest) (*api.ListImagesResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceServerMock) ImageStatus(context.Context, *api.ImageStatusRequest) (*api.ImageStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceServerMock) PullImage(context.Context, *api.PullImageRequest) (*api.PullImageResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceServerMock) RemoveImage(context.Context, *api.RemoveImageRequest) (*api.RemoveImageResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceServerMock) ImageFsInfo(context.Context, *api.ImageFsInfoRequest) (*api.ImageFsInfoResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}
