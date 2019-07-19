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
	"time"

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type RuntimeServiceProxyTestClientMock struct {
	RuntimeServiceClientMock

	t                         *testing.T
	podSandboxStatusCallCount int
	podSandboxStatusReqWant   *api.PodSandboxStatusRequest
	podSandboxStatusRes       []*api.PodSandboxStatusResponse
	podSandboxStatusErr       error
	createContainerReqWant    *api.CreateContainerRequest
	createContainerRes        *api.CreateContainerResponse
	createContainerErr        error
	startContainerReqWant     *api.StartContainerRequest
	startContainerRes         *api.StartContainerResponse
	startContainerErr         error
	containerStatusCallCount  int
	containerStatusReqWant    *api.ContainerStatusRequest
	containerStatusRes        []*api.ContainerStatusResponse
	containerStatusErr        error
}

func (m *RuntimeServiceProxyTestClientMock) PodSandboxStatus(ctx context.Context,
	req *api.PodSandboxStatusRequest, opts ...grpc.CallOption) (
	*api.PodSandboxStatusResponse, error) {
	if !reflect.DeepEqual(req, m.podSandboxStatusReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.podSandboxStatusReqWant)
	}
	var res *api.PodSandboxStatusResponse
	if l := len(m.podSandboxStatusRes); l > 0 {
		if m.podSandboxStatusCallCount >= l {
			res = m.podSandboxStatusRes[l-1]
		} else {
			res = m.podSandboxStatusRes[m.podSandboxStatusCallCount]
		}
		m.podSandboxStatusCallCount++
	}
	return res, m.podSandboxStatusErr
}

func (m *RuntimeServiceProxyTestClientMock) CreateContainer(ctx context.Context,
	req *api.CreateContainerRequest, opts ...grpc.CallOption) (
	*api.CreateContainerResponse, error) {
	if !reflect.DeepEqual(req, m.createContainerReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.createContainerReqWant)
	}
	return m.createContainerRes, m.createContainerErr
}

func (m *RuntimeServiceProxyTestClientMock) StartContainer(ctx context.Context,
	req *api.StartContainerRequest, opts ...grpc.CallOption) (*api.StartContainerResponse, error) {
	if !reflect.DeepEqual(req, m.startContainerReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.startContainerReqWant)
	}
	return m.startContainerRes, m.startContainerErr
}

func (m *RuntimeServiceProxyTestClientMock) ContainerStatus(ctx context.Context,
	req *api.ContainerStatusRequest, opts ...grpc.CallOption) (
	*api.ContainerStatusResponse, error) {
	if !reflect.DeepEqual(req, m.containerStatusReqWant) {
		m.t.Errorf("req = %s, want %s", req, m.containerStatusReqWant)
	}
	var res *api.ContainerStatusResponse
	if l := len(m.containerStatusRes); l > 0 {
		if m.containerStatusCallCount >= l {
			res = m.containerStatusRes[l-1]
		} else {
			res = m.containerStatusRes[m.containerStatusCallCount]
		}
		m.containerStatusCallCount++
	}
	return res, m.containerStatusErr
}

// TestRuntimeServiceProxy tests that at least two of the methods properly proxy calls to
// the backend client and log the calls.
func TestRuntimeServiceProxy(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &RuntimeServiceProxyTestClientMock{
		t:                       t,
		podSandboxStatusReqWant: podSandboxStatus1Req,
	}

	logger := &InMemoryLogger{}

	proxy := NewRuntimeServiceProxy(clientMock, logger)

	ctx := context.TODO()

	// Test a successful call.
	clientMock.podSandboxStatusRes = []*api.PodSandboxStatusResponse{podSandboxStatus1ResNotReady}
	clientMock.podSandboxStatusErr = nil
	gotPodSandboxStatus1Res, err := proxy.PodSandboxStatus(ctx, podSandboxStatus1Req)
	if err != nil {
		t.Errorf("PodSandboxStatus(%#v, %s) failed: %v", ctx, podSandboxStatus1Req, err)
	}
	if !reflect.DeepEqual(gotPodSandboxStatus1Res, podSandboxStatus1ResNotReady) {
		t.Errorf("PodSandboxStatus(%#v, %s) = %s, want %s", ctx, podSandboxStatus1Req,
			gotPodSandboxStatus1Res, podSandboxStatus1ResNotReady)
	}

	wantCalls := []Call{
		{
			Request:  podSandboxStatus1Req,
			Response: podSandboxStatus1ResNotReady,
			Err:      nil,
		},
	}
	if !reflect.DeepEqual(logger.Calls, wantCalls) {
		t.Errorf("Calls = %#v, want %#v", logger.Calls, wantCalls)
	}

	// Reset the logger.
	logger.Calls = nil

	// Test a call returning an error.
	clientMock.podSandboxStatusRes = nil
	clientMock.podSandboxStatusErr = customErr
	_, err = proxy.PodSandboxStatus(ctx, podSandboxStatus1Req)
	if err != customErr {
		t.Errorf("PodSandboxStatus(%#v, %s) failed with %v, want %v", ctx, podSandboxStatus1Req,
			err, customErr)
	}

	wantCalls = []Call{
		{
			Request:  podSandboxStatus1Req,
			Response: (*api.PodSandboxStatusResponse)(nil),
			Err:      customErr,
		},
	}
	if !reflect.DeepEqual(logger.Calls, wantCalls) {
		t.Errorf("Calls = %#v, want %#v", logger.Calls, wantCalls)
	}
}

func TestPodSandboxIsReadyAndHasIP(t *testing.T) {
	var tests = []struct {
		status *api.PodSandboxStatus
		want   bool
	}{
		{podSandboxStatus1ResNotReady.Status, false},
		{podSandboxStatus1ResReadyNoIP.Status, false},
		{podSandboxStatus1ResReadyIP.Status, true},
	}

	for _, test := range tests {
		got, err := PodSandboxIsReadyAndHasIP(test.status)
		if err != nil {
			t.Errorf("PodSandboxIsReadyAndHasIP(%s) failed: %v", test.status, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("PodSandboxIsReadyAndHasIP(%s) = %t, want %t", test.status, got, test.want)
		}
	}
}

func TestRuntimeServiceProxy_WaitForPodSandboxStatus(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &RuntimeServiceProxyTestClientMock{
		t:                       t,
		podSandboxStatusReqWant: podSandboxStatus1Req,
	}

	proxy := NewRuntimeServiceProxy(clientMock, &NopLogger{})

	ctx := context.TODO()

	// Fail in backend.
	clientMock.podSandboxStatusRes = nil
	clientMock.podSandboxStatusErr = customErr
	_, err := proxy.WaitForPodSandboxStatus(ctx, nil, podSandboxID1, PodSandboxIsReadyAndHasIP)
	if err == nil {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) didn't fail", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP")
	}

	// Fail in predicate.
	errorPredicate := func(status *api.PodSandboxStatus) (bool, error) {
		return false, customErr
	}
	clientMock.podSandboxStatusRes = []*api.PodSandboxStatusResponse{
		podSandboxStatus1ResReadyIP,
	}
	clientMock.podSandboxStatusErr = nil
	_, err = proxy.WaitForPodSandboxStatus(ctx, nil, podSandboxID1, errorPredicate)
	if err == nil {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) didn't fail", ctx, nil,
			podSandboxID1, "errorPredicate")
	}

	// Successful after first status check.
	clientMock.podSandboxStatusRes = []*api.PodSandboxStatusResponse{
		podSandboxStatus1ResReadyIP,
	}
	clientMock.podSandboxStatusErr = nil
	gotStatus, err := proxy.WaitForPodSandboxStatus(ctx, nil, podSandboxID1,
		PodSandboxIsReadyAndHasIP)
	if err != nil {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) failed: %v", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP", err)
	}
	if !reflect.DeepEqual(gotStatus, podSandboxStatus1ResReadyIP.Status) {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) = %s, want %s", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP", gotStatus,
			podSandboxStatus1ResReadyIP.Status)
	}

	// Successful after 3 status checks.
	clientMock.podSandboxStatusRes = []*api.PodSandboxStatusResponse{
		podSandboxStatus1ResNotReady,
		podSandboxStatus1ResReadyNoIP,
		podSandboxStatus1ResReadyIP,
	}
	clientMock.podSandboxStatusErr = nil
	gotStatus, err = proxy.WaitForPodSandboxStatus(ctx, nil, podSandboxID1,
		PodSandboxIsReadyAndHasIP)
	if err != nil {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) failed: %v", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP", err)
	}
	if !reflect.DeepEqual(gotStatus, podSandboxStatus1ResReadyIP.Status) {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) = %s, want %s", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP", gotStatus,
			podSandboxStatus1ResReadyIP.Status)
	}

	// Timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	clientMock.podSandboxStatusRes = []*api.PodSandboxStatusResponse{
		podSandboxStatus1ResNotReady,
	}
	clientMock.podSandboxStatusErr = nil
	_, err = proxy.WaitForPodSandboxStatus(ctx, nil, podSandboxID1,
		PodSandboxIsReadyAndHasIP)
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("ctx.Err() == %v, want %v", ctx.Err(), context.DeadlineExceeded)
	}
	if err != ctx.Err() {
		t.Errorf("WaitForPodSandboxStatus(%#v, %#v, %q, %s) failed with %v, want %v", ctx, nil,
			podSandboxID1, "PodSandboxIsReadyAndHasIP", err, ctx.Err())
	}
}

func TestRuntimeServiceProxy_CreateAndStartContainer(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &RuntimeServiceProxyTestClientMock{
		t:                      t,
		createContainerReqWant: createContainer1Req,
		startContainerReqWant:  startContainer1Req,
	}

	proxy := NewRuntimeServiceProxy(clientMock, &NopLogger{})

	ctx := context.TODO()

	// Fail in CreateContainer.
	clientMock.createContainerRes = nil
	clientMock.createContainerErr = customErr
	gotContainerID, err := proxy.CreateAndStartContainer(ctx, nil, podSandboxID1,
		podSandboxConfig1, createContainer1Req.Config)
	if err == nil {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) didn't fail", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config)
	}
	if gotContainerID != "" {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) = %q, want %q", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config, gotContainerID, "")
	}

	// Fail in StartContainer.
	clientMock.createContainerRes = createContainer1Res
	clientMock.createContainerErr = nil
	clientMock.startContainerRes = nil
	clientMock.startContainerErr = customErr
	gotContainerID, err = proxy.CreateAndStartContainer(ctx, nil, podSandboxID1,
		podSandboxConfig1, createContainer1Req.Config)
	if err == nil {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) didn't fail", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config)
	}
	if gotContainerID != containerID1 {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) = %q, want %q", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config, gotContainerID,
			containerID1)
	}

	// Success.
	clientMock.createContainerRes = createContainer1Res
	clientMock.createContainerErr = nil
	clientMock.startContainerRes = startContainer1Res
	clientMock.startContainerErr = nil
	gotContainerID, err = proxy.CreateAndStartContainer(ctx, nil, podSandboxID1,
		podSandboxConfig1, createContainer1Req.Config)
	if err != nil {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) failed: %v", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config, err)
	}
	if gotContainerID != containerID1 {
		t.Errorf("CreateAndStartContainer(%#v, %#v, %q, %s, %s) = %q, want %q", ctx, nil,
			podSandboxID1, podSandboxConfig1, createContainer1Req.Config, gotContainerID,
			containerID1)
	}
}

func TestRuntimeServiceProxy_WaitForContainerStatus(t *testing.T) {
	customErr := errors.New("testing error")

	clientMock := &RuntimeServiceProxyTestClientMock{
		t:                      t,
		containerStatusReqWant: containerStatus1Req,
	}

	proxy := NewRuntimeServiceProxy(clientMock, &NopLogger{})

	ctx := context.TODO()

	containerExited := func(status *api.ContainerStatus) (bool, error) {
		return status.GetState() == api.ContainerState_CONTAINER_EXITED, nil
	}

	// Fail in backend.
	clientMock.containerStatusRes = nil
	clientMock.containerStatusErr = customErr
	_, err := proxy.WaitForContainerStatus(ctx, nil, containerID1, containerExited)
	if err == nil {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) didn't fail", ctx, nil, containerID1,
			"containerExited")
	}

	// Fail in predicate.
	errorPredicate := func(status *api.ContainerStatus) (bool, error) {
		return false, customErr
	}
	clientMock.containerStatusRes = []*api.ContainerStatusResponse{
		containerStatus1ResCreated,
	}
	clientMock.containerStatusErr = nil
	_, err = proxy.WaitForContainerStatus(ctx, nil, containerID1, errorPredicate)
	if err == nil {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) didn't fail", ctx, nil, containerID1,
			"errorPredicate")
	}

	// Successful after first status check.
	clientMock.containerStatusRes = []*api.ContainerStatusResponse{
		containerStatus1ResExited,
	}
	clientMock.containerStatusErr = nil
	gotStatus, err := proxy.WaitForContainerStatus(ctx, nil, containerID1, containerExited)
	if err != nil {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) failed: %v", ctx, nil, containerID1,
			"containerExited", err)
	}
	if !reflect.DeepEqual(gotStatus, containerStatus1ResExited.Status) {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) = %s, want %s", ctx, nil, containerID1,
			"containerExited", gotStatus, containerStatus1ResExited.Status)
	}

	// Successful after 3 status checks.
	clientMock.containerStatusRes = []*api.ContainerStatusResponse{
		containerStatus1ResCreated,
		containerStatus1ResRunning,
		containerStatus1ResExited,
	}
	clientMock.containerStatusErr = nil
	gotStatus, err = proxy.WaitForContainerStatus(ctx, nil, containerID1, containerExited)
	if err != nil {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) failed: %v", ctx, nil, containerID1,
			"containerExited", err)
	}
	if !reflect.DeepEqual(gotStatus, containerStatus1ResExited.Status) {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) = %s, want %s", ctx, nil, containerID1,
			"containerExited", gotStatus, containerStatus1ResExited.Status)
	}

	// Timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	clientMock.containerStatusRes = []*api.ContainerStatusResponse{
		containerStatus1ResCreated,
	}
	clientMock.containerStatusErr = nil
	_, err = proxy.WaitForContainerStatus(ctx, nil, containerID1, containerExited)
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("ctx.Err() == %v, want %v", ctx.Err(), context.DeadlineExceeded)
	}
	if err != ctx.Err() {
		t.Errorf("WaitForContainerStatus(%#v, %#v, %q, %s) failed with %v, want %v", ctx, nil,
			containerID1, "containerExited", err, ctx.Err())
	}
}
