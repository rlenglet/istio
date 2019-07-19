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
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type ClientTestRuntimeServiceServerMock struct {
	RuntimeServiceServerMock

	t    *testing.T
	want *api.VersionRequest
	res  *api.VersionResponse
}

func (m *ClientTestRuntimeServiceServerMock) Version(ctx context.Context, req *api.VersionRequest) (
	*api.VersionResponse, error) {
	if !reflect.DeepEqual(req, m.want) {
		m.t.Errorf("req = %s, want %s", req, m.want)
	}
	return m.res, nil
}

type ClientTestImageServiceServerMock struct {
	ImageServiceServerMock

	t    *testing.T
	want *api.ListImagesRequest
	res  *api.ListImagesResponse
}

func (m *ClientTestImageServiceServerMock) ListImages(ctx context.Context, req *api.ListImagesRequest) (
	*api.ListImagesResponse, error) {
	if !reflect.DeepEqual(req, m.want) {
		m.t.Errorf("req = %s, want %s", req, m.want)
	}
	return m.res, nil
}

func TestNewClientsAndServer(t *testing.T) {
	// Start mock servers listening on different Unix domain sockets with random names.
	const dir = ""
	const pattern = "*.sock"
	socket1, err := ioutil.TempFile("", "*.sock")
	if err != nil {
		t.Fatalf("ioutil.TempFile(%q, %q) failed: %v", dir, pattern, err)
	}
	socket2, err := ioutil.TempFile("", "*.sock")
	if err != nil {
		t.Fatalf("ioutil.TempFile(%q, %q) failed: %v", dir, pattern, err)
	}
	socket1Name := socket1.Name()
	socket2Name := socket2.Name()
	if err := os.Remove(socket1Name); err != nil {
		t.Fatalf("os.Remove(%q) failed: %v", socket1Name, err)
	}
	if err := os.Remove(socket2Name); err != nil {
		t.Fatalf("os.Remove(%q) failed: %v", socket2Name, err)
	}

	server1 := NewServer(socket1Name,
		&ClientTestRuntimeServiceServerMock{
			t:    t,
			want: version1Req,
			res:  version1Res,
		},
		&ClientTestImageServiceServerMock{
			t:    t,
			want: listImages1Req,
			res:  listImages1Res,
		})
	server2 := NewServer(socket2Name,
		&ClientTestRuntimeServiceServerMock{
			t:    t,
			want: version2Req,
			res:  version2Res,
		},
		&ClientTestImageServiceServerMock{
			t:    t,
			want: listImages2Req,
			res:  listImages2Res,
		})

	defer func() {
		server1.Stop()
		server2.Stop()
		_ = os.Remove(socket1Name)
		_ = os.Remove(socket2Name)
	}()

	if err := server1.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if err := server2.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	ctx := context.TODO()

	// Test clients pointing to the same server.
	runtimeServiceClient1, imageServiceClient1, err := NewClients(context.TODO(),
		socket1Name, socket1Name)
	if err != nil {
		t.Fatalf("NewClients(%#v, %q, %q) failed: %v", ctx, socket1Name, socket1Name, err)
	}

	gotVersion1Res, err := runtimeServiceClient1.Version(ctx, version1Req)
	if err != nil {
		t.Errorf("Version(%#v, %s) failed: %v", ctx, version1Req, err)
	}
	if !reflect.DeepEqual(gotVersion1Res, version1Res) {
		t.Errorf("Version(%#v, %s) = %s, want %s", ctx, version1Req,
			gotVersion1Res, version1Res)
	}

	gotListImages1Res, err := imageServiceClient1.ListImages(ctx, listImages1Req)
	if err != nil {
		t.Errorf("ListImages(%#v, %s) failed: %v", ctx, listImages1Req, err)
	}
	if !reflect.DeepEqual(gotListImages1Res, listImages1Res) {
		t.Errorf("ListImages(%#v, %s) = %s, want %s", ctx, listImages1Req,
			gotListImages1Res, listImages1Res)
	}

	// Test clients pointing to different servers.
	runtimeServiceClient2, imageServiceClient2, err := NewClients(context.TODO(),
		socket1Name, socket2Name)
	if err != nil {
		t.Fatalf("NewClients(%#v, %q, %q) failed: %v", ctx, socket1Name, socket2Name, err)
	}

	gotVersion1Res, err = runtimeServiceClient2.Version(ctx, version1Req)
	if err != nil {
		t.Errorf("Version(%#v, %s) failed: %v", ctx, version1Req, err)
	}
	if !reflect.DeepEqual(gotVersion1Res, version1Res) {
		t.Errorf("Version(%#v, %s) = %s, want %s", ctx, version1Req,
			gotVersion1Res, version1Res)
	}

	gotListImages2Res, err := imageServiceClient2.ListImages(ctx, listImages2Req)
	if err != nil {
		t.Errorf("ListImages(%#v, %s) failed: %v", ctx, listImages2Req, err)
	}
	if !reflect.DeepEqual(gotListImages2Res, listImages2Res) {
		t.Errorf("ListImages(%#v, %s) = %s, want %s", ctx, listImages2Req,
			gotListImages2Res, listImages2Res)
	}

	// Test error conditions.
	runtimeServiceClient3, _, err := NewClients(context.TODO(), invalidSocketName, socket2Name)
	if err != nil {
		t.Fatalf("NewClients(%#v, %q, %q) failed: %v", ctx, invalidSocketName, socket2Name, err)
	} else {
		_, err = runtimeServiceClient3.Version(ctx, version1Req)
		if err == nil {
			t.Errorf("Version(%#v, %s) didn't fail", ctx, version1Req)
		}
		if status.Code(err) != codes.Unavailable {
			t.Errorf("Version(%#v, %s) failed with status code %q, want %q", ctx,
				version1Req, status.Code(err), codes.Unavailable)
		}
	}

	_, imageServiceClient3, err := NewClients(context.TODO(), socket1Name, invalidSocketName)
	if err != nil {
		t.Errorf("NewClients(%#v, %q, %q) failed: %v", ctx, socket1Name, invalidSocketName, err)
	} else {
		_, err = imageServiceClient3.ListImages(ctx, listImages1Req)
		if err == nil {
			t.Errorf("ListImages(%#v, %s) didn't fail", ctx, listImages1Req)
		}
		if status.Code(err) != codes.Unavailable {
			t.Errorf("ListImages(%#v, %s) failed with status code %q, want %q", ctx,
				listImages1Req, status.Code(err), codes.Unavailable)
		}
	}
}
