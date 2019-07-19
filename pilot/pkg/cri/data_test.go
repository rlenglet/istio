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
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	invalidSocketName = "invalid-socket-name"

	podSandboxID1 = "pod-sandbox1"
	containerID1  = "container1"
	image1        = "image1"
	imageRef1     = "image-ref1"
	image2        = "image2"
	imageRef2     = "image-ref2"
)

var (
	version1Req = &api.VersionRequest{Version: "1.0.1"}
	version1Res = &api.VersionResponse{Version: "1.0.2"}
	version2Req = &api.VersionRequest{Version: "2.0.1"}
	version2Res = &api.VersionResponse{Version: "2.0.2"}

	podSandboxStatus1Req         = &api.PodSandboxStatusRequest{PodSandboxId: podSandboxID1}
	podSandboxStatus1ResNotReady = &api.PodSandboxStatusResponse{Status: &api.PodSandboxStatus{
		Id:    podSandboxID1,
		State: api.PodSandboxState_SANDBOX_NOTREADY,
	}}
	podSandboxStatus1ResReadyNoIP = &api.PodSandboxStatusResponse{Status: &api.PodSandboxStatus{
		Id:    podSandboxID1,
		State: api.PodSandboxState_SANDBOX_READY,
	}}
	podSandboxStatus1ResReadyIP = &api.PodSandboxStatusResponse{Status: &api.PodSandboxStatus{
		Id:    podSandboxID1,
		State: api.PodSandboxState_SANDBOX_READY,
		Network: &api.PodSandboxNetworkStatus{
			Ip: "10.0.0.123",
		},
	}}

	createContainer1Req = &api.CreateContainerRequest{
		PodSandboxId: podSandboxID1,
		Config: &api.ContainerConfig{
			Metadata: &api.ContainerMetadata{
				Name: "helloworld",
			},
			Image: &api.ImageSpec{
				Image: imageRef1,
			},
			// TODO: Complete.
		},
		SandboxConfig: podSandboxConfig1,
	}
	createContainer1Res = &api.CreateContainerResponse{
		ContainerId: containerID1,
	}

	startContainer1Req = &api.StartContainerRequest{
		ContainerId: containerID1,
	}
	startContainer1Res = &api.StartContainerResponse{}

	containerStatus1Req        = &api.ContainerStatusRequest{ContainerId: containerID1}
	containerStatus1ResCreated = &api.ContainerStatusResponse{Status: &api.ContainerStatus{
		Id:    containerID1,
		State: api.ContainerState_CONTAINER_CREATED,
	}}
	containerStatus1ResRunning = &api.ContainerStatusResponse{Status: &api.ContainerStatus{
		Id:    containerID1,
		State: api.ContainerState_CONTAINER_RUNNING,
	}}
	containerStatus1ResExited = &api.ContainerStatusResponse{Status: &api.ContainerStatus{
		Id:    containerID1,
		State: api.ContainerState_CONTAINER_EXITED,
	}}

	listImages1Req = &api.ListImagesRequest{
		Filter: &api.ImageFilter{Image: &api.ImageSpec{Image: image1}},
	}
	listImages1Res = &api.ListImagesResponse{
		Images: []*api.Image{{Id: imageRef1}},
	}
	listImages2Req = &api.ListImagesRequest{
		Filter: &api.ImageFilter{Image: &api.ImageSpec{Image: image2}},
	}
	listImages2Res = &api.ListImagesResponse{
		Images: []*api.Image{{Id: imageRef2}},
	}

	authConfig1 = &api.AuthConfig{
		Username: "alice",
		Password: "password1!",
	}
	podSandboxConfig1 = &api.PodSandboxConfig{
		Metadata: &api.PodSandboxMetadata{
			Name:      "helloworld",
			Namespace: "default",
		},
	}
	pullImage1Req = &api.PullImageRequest{
		Image: &api.ImageSpec{
			Image: image1,
		},
		Auth:          authConfig1,
		SandboxConfig: podSandboxConfig1,
	}
	pullImage1Res = &api.PullImageResponse{
		ImageRef: imageRef1,
	}
)
