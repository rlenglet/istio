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

// ExtendedImageServiceServer extends ImageServiceServer with higher-level utility functions.
type ExtendedImageServiceServer interface {
	api.ImageServiceServer

	// ImageRef returns the image reference for the given image, first pulling if it is not yet
	// pulled.
	ImageRef(ctx context.Context, logFields []zap.Field, podSandboxConfig *api.PodSandboxConfig,
		image string, authConfig *api.AuthConfig) (string, error)
}

// ImageServiceProxy is an ImageServiceServer that forwards all received calls to a backend
// ImageServiceClient. This is intended to be overridden to intercept specific calls.
type ImageServiceProxy struct {
	// imageService is the backend ImageServiceClient that calls are forwarded to.
	imageService api.ImageServiceClient

	// runtimeServiceLogger logs calls forwarded to imageService.
	imageServiceLogger Logger
}

var _ ExtendedImageServiceServer = &ImageServiceProxy{}

// NewImageServiceProxy creates a new ImageServiceProxy backed by the given ImageServiceClient.
func NewImageServiceProxy(imageService api.ImageServiceClient, imageServiceLogger Logger) *ImageServiceProxy {
	return &ImageServiceProxy{
		imageService:       imageService,
		imageServiceLogger: imageServiceLogger,
	}
}

func (p *ImageServiceProxy) ImageRef(ctx context.Context, logFields []zap.Field,
	podSandboxConfig *api.PodSandboxConfig, image string, authConfig *api.AuthConfig) (
	string, error) {
	imageSpec := api.ImageSpec{
		Image: image,
	}

	var imageRef string

	listReq := &api.ListImagesRequest{
		Filter: &api.ImageFilter{
			Image: &imageSpec,
		},
	}
	listRes, err := p.ListImages(ctx, listReq)
	if err != nil {
		return "", err
	}
	for _, img := range listRes.GetImages() {
		for _, tag := range img.GetRepoTags() {
			if tag == image {
				log.Debug("Image was already pulled", logFields...)
				return img.Id, nil
			}
		}
	}

	log.Info("Pulling image", logFields...)

	pullReq := &api.PullImageRequest{
		Image:         &imageSpec,
		Auth:          authConfig,
		SandboxConfig: podSandboxConfig,
	}
	pullRes, err := p.PullImage(ctx, pullReq)
	if err != nil {
		return "", err
	}
	imageRef = pullRes.GetImageRef()

	return imageRef, nil
}

func (p *ImageServiceProxy) ListImages(ctx context.Context, req *api.ListImagesRequest) (*api.ListImagesResponse, error) {
	res, err := p.imageService.ListImages(ctx, req)
	p.imageServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *ImageServiceProxy) ImageStatus(ctx context.Context, req *api.ImageStatusRequest) (*api.ImageStatusResponse, error) {
	res, err := p.imageService.ImageStatus(ctx, req)
	p.imageServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *ImageServiceProxy) PullImage(ctx context.Context, req *api.PullImageRequest) (*api.PullImageResponse, error) {
	res, err := p.imageService.PullImage(ctx, req)
	p.imageServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *ImageServiceProxy) RemoveImage(ctx context.Context, req *api.RemoveImageRequest) (*api.RemoveImageResponse, error) {
	res, err := p.imageService.RemoveImage(ctx, req)
	p.imageServiceLogger.LogCall(req, res, err)
	return res, err
}

func (p *ImageServiceProxy) ImageFsInfo(ctx context.Context, req *api.ImageFsInfoRequest) (*api.ImageFsInfoResponse, error) {
	res, err := p.imageService.ImageFsInfo(ctx, req)
	p.imageServiceLogger.LogCall(req, res, err)
	return res, err
}
