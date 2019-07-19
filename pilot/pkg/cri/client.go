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
	"go.uber.org/zap"
	"net"

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// NewClients creates new RuntimeServiceClient and ImageServiceClient connected to CRI server(s)
// via the given Unix domain socket paths.
// If the socket paths are the same, the same connection is used for both clients.
func NewClients(ctx context.Context, runtimeServiceSocketPath, imageServiceSocketPath string) (
	api.RuntimeServiceClient, api.ImageServiceClient, error) {
	cc, err := dial(ctx, runtimeServiceSocketPath)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to connect to RuntimeService CRI server using socket file %q: %v",
			runtimeServiceSocketPath, err)
	}
	runtimeServiceClient := api.NewRuntimeServiceClient(cc)

	// If the socket paths are the same, reuse the gRPC connection.
	if runtimeServiceSocketPath != imageServiceSocketPath {
		cc, err = dial(ctx, imageServiceSocketPath)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to connect to ImageService CRI server using socket file %q: %v",
				imageServiceSocketPath, err)
		}
	}
	imageServiceClient := api.NewImageServiceClient(cc)

	return runtimeServiceClient, imageServiceClient, nil
}

// dial connects to a CRI server via the given Unix domain socket path.
func dial(ctx context.Context, socketPath string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, socketPath, grpc.WithInsecure(),
		grpc.WithContextDialer(contextDialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
}

func contextDialer(ctx context.Context, path string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", path)
	if err != nil {
		log.Warn("Failed to connect to CRI server",
			zap.String("endpoint", path),
			zap.Error(err))
	}
	return conn, err
}
