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

	istiolog "istio.io/pkg/log"
)

var log = istiolog.RegisterScope("cri", "cri", 0)

// WithField appends a field into the array of log fields.
func WithField(newField zap.Field, logFields ...zap.Field) []zap.Field {
	if len(logFields) == 0 {
		return []zap.Field{newField}
	}

	newFields := make([]zap.Field, 0, len(logFields)+1)
	for _, f := range logFields {
		newFields = append(newFields, f)
	}
	newFields = append(newFields, newField)
	return newFields
}

// WithError appends an error field into the array of log fields.
func WithError(err error, logFields ...zap.Field) []zap.Field {
	return WithField(zap.Error(err), logFields...)
}
