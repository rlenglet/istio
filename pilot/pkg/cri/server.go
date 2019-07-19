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
	"fmt"
	"net"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	// maxMsgSize is the default message size limit (16MB).
	maxMsgSize = 1024 * 1024 * 16
)

// Server is a gRPC server implementing CRI's RuntimeService and ImageService over a Unix domain
// socket.
type Server struct {
	// socketPath is the absolute path to the Unix domain socket where the server is listening.
	socketPath string

	// runtimeService is the server implementation of CRI's RuntimeService.
	runtimeService api.RuntimeServiceServer

	// imageService is the server implementation of CRI's ImageService.
	imageService api.ImageServiceServer

	// server is the gRPC server.
	server *grpc.Server
}

// NewServer creates a new server to serve the given services. The socketPath must be absolute.
func NewServer(socketPath string, runtimeService api.RuntimeServiceServer,
	imageService api.ImageServiceServer) *Server {
	return &Server{
		socketPath:     socketPath,
		runtimeService: runtimeService,
		imageService:   imageService,
	}
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	log.Info("Starting CRI server", zap.String("server_endpoint", s.socketPath))

	l, err := CreateSocketListener(s.socketPath)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	api.RegisterRuntimeServiceServer(s.server, s.runtimeService)
	api.RegisterImageServiceServer(s.server, s.imageService)

	go func() {
		if err := s.server.Serve(l); err != nil {
			log.Fatal("Failed to serve CRI connections",
				zap.String("server_endpoint", s.socketPath),
				zap.Error(err))
		}
	}()
	return nil
}

// Stop stops the gRPC server.
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

// CreateSocketListener creates a Unix domain socket listener.
func CreateSocketListener(socketPath string) (net.Listener, error) {
	dir := filepath.Dir(socketPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create socket directory %q: %v", dir, err)
	}

	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove previous socket file %q: %v", socketPath, err)
	}

	return net.Listen("unix", socketPath)
}
