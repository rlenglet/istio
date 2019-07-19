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

package kubeinject

import (
	"fmt"

	k8s "k8s.io/api/core/v1"

	"istio.io/istio/pkg/kube/inject"
)

const (
	// ContainerNameIstioInit is the required name for an injected init container.
	ContainerNameIstioInit = "istio-init"

	// ContainerNameIstioProxy is the required name for an injected proxy container.
	ContainerNameIstioProxy = "istio-proxy"
)

// ValidateSidecarInjectionSpec validates that the given spec can be translated.
func ValidateSidecarInjectionSpec(spec *inject.SidecarInjectionSpec) error {
	// TODO(rlenglet): Check whether spec.RewriteAppHTTPProbe needs to be handled specially.

	if l := len(spec.InitContainers); l != 1 {
		return fmt.Errorf("spec has %d init containers, expected 1", l)
	}
	initContainer := &spec.InitContainers[0]
	if initContainer.Name != ContainerNameIstioInit {
		return fmt.Errorf("init container has name %q, expected %q", initContainer.Name,
			ContainerNameIstioInit)
	}
	if l := len(initContainer.Ports); l != 0 {
		return fmt.Errorf("init container has %d ports, expected 0", l)
	}
	if err := ValidateContainer(initContainer); err != nil {
		return fmt.Errorf("init container is invalid: %v", err)
	}

	if l := len(spec.Containers); l != 1 {
		return fmt.Errorf("spec has %d containers, expected 1", l)
	}
	proxyContainer := &spec.Containers[0]
	if proxyContainer.Name != ContainerNameIstioProxy {
		return fmt.Errorf("proxy container has name %q, expected %q", proxyContainer.Name,
			ContainerNameIstioProxy)
	}
	if err := ValidateContainer(proxyContainer); err != nil {
		return fmt.Errorf("proxy container is invalid: %v", err)
	}

	// TODO(rlenglet): spec.Volumes

	if err := ValidatePodDNSConfig(spec.DNSConfig); err != nil {
		return err
	}

	// TODO(rlenglet): spec.ImagePullSecrets

	return nil
}

// ValidateContainer validates that the given container spec can be translated.
func ValidateContainer(container *k8s.Container) error {
	if l := len(container.EnvFrom); l != 0 {
		return fmt.Errorf("container has %d envFrom sources, expected 0", l)
	}
	if err := ValidateEnv(container.Env); err != nil {
		return err
	}
	if err := ValidateResourceList(container.Resources.Limits); err != nil {
		return fmt.Errorf("unsupported resource limits: %v", err)
	}
	if err := ValidateResourceList(container.Resources.Requests); err != nil {
		return fmt.Errorf("unsupported resource requests: %v", err)
	}

	// TODO(rlenglet): Volume mounts, etc.

	return nil
}

// ValidateEnv validates that the given environment variables can be translated.
func ValidateEnv(envVars []k8s.EnvVar) error {
	for _, env := range envVars {
		if env.Value != "" && env.ValueFrom != nil {
			return fmt.Errorf("env %q has both value and valueFrom", env.Name)
		}
		if env.ValueFrom != nil {
			if env.ValueFrom.FieldRef == nil ||
				env.ValueFrom.ResourceFieldRef != nil ||
				env.ValueFrom.ConfigMapKeyRef != nil ||
				env.ValueFrom.SecretKeyRef != nil {
				return fmt.Errorf("env %q has neither value nor fieldRef", env.Name)
			}
			switch env.ValueFrom.FieldRef.FieldPath {
			case "metadata.name":
			case "metadata.namespace":
			case "status.hostIP":
			case "status.podIP":
			default:
				return fmt.Errorf("env %q uses unsupported fieldPath", env.Name)
			}
		}
	}

	return nil
}

// ValidateResourceList validates that the given resource list can be translated.
func ValidateResourceList(resources k8s.ResourceList) error {
	for resourceName := range resources {
		switch resourceName {
		case k8s.ResourceCPU:
		case k8s.ResourceMemory:
		default:
			return fmt.Errorf("unsupported resource %q", resourceName)
		}
	}

	return nil
}

// ValidatePodDNSConfig validates that the given pod DNS config can be translated.
func ValidatePodDNSConfig(config *k8s.PodDNSConfig) error {
	if config == nil {
		return nil
	}

	if l := len(config.Nameservers); l != 0 {
		return fmt.Errorf("DNS config has %d nameservers, expected 0", l)
	}
	if l := len(config.Options); l != 0 {
		return fmt.Errorf("DNS config has %d options, expected 0", l)
	}

	return nil
}
