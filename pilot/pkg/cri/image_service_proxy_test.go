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
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type ImageServiceProxyTestClientMock struct {
	ImageServiceClientMock

	t                 *testing.T
	listImagesReqWant *api.ListImagesRequest
	listImagesRes     *api.ListImagesResponse
	listImagesErr     error
	pullImageReqWant  *api.PullImageRequest
	pullImageRes      *api.PullImageResponse
	pullImageErr      error
}

func (m *ImageServiceProxyTestClientMock) ListImages(ctx context.Context,
	req *api.ListImagesRequest, opts ...grpc.CallOption) (*api.ListImagesResponse, error) {
	if !reflect.DeepEqual(req, m.listImagesReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.listImagesReqWant)
	}
	return m.listImagesRes, m.listImagesErr
}

func (m *ImageServiceProxyTestClientMock) PullImage(ctx context.Context, req *api.PullImageRequest,
	opts ...grpc.CallOption) (*api.PullImageResponse, error) {
	if !reflect.DeepEqual(req, m.pullImageReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.pullImageReqWant)
	}
	return m.pullImageRes, m.pullImageErr
}

// TestImageServiceProxy tests that at least two of the methods properly proxy calls to
// the backend client and log the calls.
func TestImageServiceProxy(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &ImageServiceProxyTestClientMock{
		t:                 t,
		listImagesReqWant: listImages1Req,
	}

	logger := &InMemoryLogger{}

	proxy := NewImageServiceProxy(clientMock, logger)

	ctx := context.TODO()

	// Test a successful call.
	clientMock.listImagesRes = listImages1Res
	clientMock.listImagesErr = nil
	gotListImages1Res, err := proxy.ListImages(ctx, listImages1Req)
	if err != nil {
		t.Errorf("ListImages(%#v, %s) failed: %v", ctx, listImages1Req, err)
	}
	if !reflect.DeepEqual(gotListImages1Res, listImages1Res) {
		t.Errorf("ListImages(%#v, %s) = %s, want %s", ctx, listImages1Req,
			gotListImages1Res, listImages1Res)
	}

	wantCalls := []Call{
		{
			Request:  listImages1Req,
			Response: listImages1Res,
			Err:      nil,
		},
	}
	if !reflect.DeepEqual(logger.Calls, wantCalls) {
		t.Errorf("Calls = %#v, want %#v", logger.Calls, wantCalls)
	}

	// Reset the logger.
	logger.Calls = nil

	// Test a call returning an error.
	clientMock.listImagesRes = nil
	clientMock.listImagesErr = customErr
	_, err = proxy.ListImages(ctx, listImages1Req)
	if err != customErr {
		t.Errorf("ListImages(%#v, %s) failed with %v, want %v", ctx, listImages1Req, err,
			customErr)
	}

	wantCalls = []Call{
		{
			Request:  listImages1Req,
			Response: (*api.ListImagesResponse)(nil),
			Err:      customErr,
		},
	}
	if !reflect.DeepEqual(logger.Calls, wantCalls) {
		t.Errorf("Calls = %#v, want %#v", logger.Calls, wantCalls)
	}
}

func TestImageServiceProxy_ImageRef(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &ImageServiceProxyTestClientMock{
		t:                 t,
		listImagesReqWant: listImages1Req,
		pullImageReqWant:  pullImage1Req,
	}

	proxy := NewImageServiceProxy(clientMock, &NopLogger{})

	ctx := context.TODO()

	clientMock.listImagesRes = nil
	clientMock.listImagesErr = customErr
	clientMock.pullImageRes = nil
	clientMock.pullImageErr = nil
	_, err := proxy.ImageRef(ctx, nil, podSandboxConfig1, image1, authConfig1)
	if err != customErr {
		t.Errorf("ImageRef(%#v, %#v, %q, %q, %s) failed with %v, want %v", ctx, nil,
			podSandboxConfig1, image1, authConfig1, err, customErr)
	}

	clientMock.listImagesRes = listImages1Res
	clientMock.listImagesErr = nil
	clientMock.pullImageRes = nil
	clientMock.pullImageErr = customErr
	_, err = proxy.ImageRef(ctx, nil, podSandboxConfig1, image1, authConfig1)
	if err != customErr {
		t.Errorf("ImageRef(%#v, %#v, %q, %q, %s) failed with %v, want %v", ctx, nil,
			podSandboxConfig1, image1, authConfig1, err, customErr)
	}

	clientMock.listImagesRes = listImages1Res
	clientMock.listImagesErr = nil
	clientMock.pullImageRes = pullImage1Res
	clientMock.pullImageErr = nil
	gotImageRef, err := proxy.ImageRef(ctx, nil, podSandboxConfig1, image1, authConfig1)
	if err != nil {
		t.Errorf("ImageRef(%#v, %#v, %q, %q, %s) failed: %v", ctx, nil, podSandboxConfig1, image1,
			authConfig1, err)
	}
	if gotImageRef != imageRef1 {
		t.Errorf("ImageRef(%#v, %#v, %q, %q, %s) = %q, want %q", ctx, nil, podSandboxConfig1,
			image1, authConfig1, gotImageRef, imageRef1)
	}
}
