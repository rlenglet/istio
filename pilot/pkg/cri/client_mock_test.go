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

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RuntimeServiceClientMock is a super type of mock implementations of RuntimeServiceClient for
// testing purposes.
// This mock's methods should be overridden to return proper responses instead of returning errors.
type RuntimeServiceClientMock struct{}

var _ api.RuntimeServiceClient = &RuntimeServiceClientMock{}

func (m *RuntimeServiceClientMock) Version(context.Context, *api.VersionRequest, ...grpc.CallOption) (*api.VersionResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) RunPodSandbox(context.Context, *api.RunPodSandboxRequest, ...grpc.CallOption) (*api.RunPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) StopPodSandbox(context.Context, *api.StopPodSandboxRequest, ...grpc.CallOption) (*api.StopPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) RemovePodSandbox(context.Context, *api.RemovePodSandboxRequest, ...grpc.CallOption) (*api.RemovePodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) PodSandboxStatus(context.Context, *api.PodSandboxStatusRequest, ...grpc.CallOption) (*api.PodSandboxStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ListPodSandbox(context.Context, *api.ListPodSandboxRequest, ...grpc.CallOption) (*api.ListPodSandboxResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) CreateContainer(context.Context, *api.CreateContainerRequest, ...grpc.CallOption) (*api.CreateContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) StartContainer(context.Context, *api.StartContainerRequest, ...grpc.CallOption) (*api.StartContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) StopContainer(context.Context, *api.StopContainerRequest, ...grpc.CallOption) (*api.StopContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) RemoveContainer(context.Context, *api.RemoveContainerRequest, ...grpc.CallOption) (*api.RemoveContainerResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ListContainers(context.Context, *api.ListContainersRequest, ...grpc.CallOption) (*api.ListContainersResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ContainerStatus(context.Context, *api.ContainerStatusRequest, ...grpc.CallOption) (*api.ContainerStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) UpdateContainerResources(context.Context, *api.UpdateContainerResourcesRequest, ...grpc.CallOption) (*api.UpdateContainerResourcesResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ReopenContainerLog(context.Context, *api.ReopenContainerLogRequest, ...grpc.CallOption) (*api.ReopenContainerLogResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ExecSync(context.Context, *api.ExecSyncRequest, ...grpc.CallOption) (*api.ExecSyncResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) Exec(context.Context, *api.ExecRequest, ...grpc.CallOption) (*api.ExecResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) Attach(context.Context, *api.AttachRequest, ...grpc.CallOption) (*api.AttachResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) PortForward(context.Context, *api.PortForwardRequest, ...grpc.CallOption) (*api.PortForwardResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ContainerStats(context.Context, *api.ContainerStatsRequest, ...grpc.CallOption) (*api.ContainerStatsResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) ListContainerStats(context.Context, *api.ListContainerStatsRequest, ...grpc.CallOption) (*api.ListContainerStatsResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) UpdateRuntimeConfig(context.Context, *api.UpdateRuntimeConfigRequest, ...grpc.CallOption) (*api.UpdateRuntimeConfigResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *RuntimeServiceClientMock) Status(context.Context, *api.StatusRequest, ...grpc.CallOption) (*api.StatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

// ImageServiceClientMock is a super type of mock implementations of ImageServiceClient for
// testing purposes.
// This mock's methods should be overridden to return proper responses instead of returning errors.
type ImageServiceClientMock struct{}

var _ api.ImageServiceClient = &ImageServiceClientMock{}

func (m *ImageServiceClientMock) ListImages(context.Context, *api.ListImagesRequest, ...grpc.CallOption) (*api.ListImagesResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceClientMock) ImageStatus(context.Context, *api.ImageStatusRequest, ...grpc.CallOption) (*api.ImageStatusResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceClientMock) PullImage(context.Context, *api.PullImageRequest, ...grpc.CallOption) (*api.PullImageResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceClientMock) RemoveImage(context.Context, *api.RemoveImageRequest, ...grpc.CallOption) (*api.RemoveImageResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}

func (m *ImageServiceClientMock) ImageFsInfo(context.Context, *api.ImageFsInfoRequest, ...grpc.CallOption) (*api.ImageFsInfoResponse, error) {
	return nil, fmt.Errorf("mock method was not overridden")
}
