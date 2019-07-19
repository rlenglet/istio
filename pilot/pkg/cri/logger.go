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
	"go.uber.org/zap"
)

// Logger logs the request, response, and error of a CRI call.
type Logger interface {
	LogCall(request, response interface{}, err error)
}

// logger is a Logger that logs all RPCs at Debug level.
type logger struct{}

var _ Logger = &logger{}

func (l *logger) LogCall(request, response interface{}, err error) {
	if err != nil {
		log.Debug("CRI RPC failed", zap.Any("request", request), zap.Any("response", response),
			zap.Error(err))
	} else {
		log.Debug("CRI RPC succeeded", zap.Any("request", request), zap.Any("response", response))
	}
}

// DefaultLogger is a Logger that logs all RPCs at Debug level.
var DefaultLogger Logger = &logger{}
