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
	"context"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	k8s "k8s.io/api/core/v1"
	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	"istio.io/istio/pilot/pkg/cri"
	"istio.io/istio/pkg/kube/inject"
)

const (
	// cpuSharesMin is the minimum number of CPU shares allocated to any container.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go#L78-L81
	cpuSharesMin = 2

	// cpuSharesPerCPU is the number of CPU shares corresponding to 100% of one CPU.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go#L40
	cpuSharesPerCPU = 1024

	// milliCPUsPerCPU is the number of milli-CPUs per CPU.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go#L41
	milliCPUsPerCPU = 1000

	// cpuPeriodDefault is the default CPU CFS (Completely Fair Scheduler) period, in milliseconds.
	// 100000 is equivalent to 100ms.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go#L44
	cpuPeriodDefault = 100000

	// cpuQuotaMin is the minimum CPU CFS (Completely Fair Scheduler) quota (1 millisecond).
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go#L45
	cpuQuotaMin = 1000

	// oomScoreAdjBestEffort is the amount by which the OOM score of best-effort processes in a
	// container should be adjusted.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/6796645672b4de8d2b6fae780c660926a23e326a/pkg/kubelet/qos/policy.go#L35
	oomScoreAdjBestEffort = 1000

	// ApparmorContainerAnnotationKeyPrefix is the prefix to an annotation key specifying a
	// container's Apparmor profile.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/security/apparmor/helpers.go#L28
	ApparmorContainerAnnotationKeyPrefix = "container.apparmor.security.beta.kubernetes.io/"
)

var (
	// defaultMaskedPaths is the default set of paths that should be masked by the container
	// runtime.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/securitycontext/util.go#L157-L167
	defaultMaskedPaths = []string{
		"/proc/acpi",
		"/proc/kcore",
		"/proc/keys",
		"/proc/latency_stats",
		"/proc/timer_list",
		"/proc/timer_stats",
		"/proc/sched_debug",
		"/proc/scsi",
		"/sys/firmware",
	}

	// defaultReadonlyPaths is the default set of paths that should be set as readonly by the
	// container runtime.
	//
	// Cf. https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/securitycontext/util.go#L168-L175
	defaultReadonlyPaths = []string{
		"/proc/asound",
		"/proc/bus",
		"/proc/fs",
		"/proc/irq",
		"/proc/sys",
		"/proc/sysrq-trigger",
	}
)

// Config is the set of parameters for injecting pods on this host.
type Config struct {
	// HostIP is the IP address of this host.
	HostIP string

	// SeccompProfileRoot is the directory path for seccomp profiles.
	//
	// Cf. kubelet's flag:
	// https://github.com/kubernetes/kubernetes/blob/c7f9dd0bafe2f3328148a85eed9c18322b9f308e/cmd/kubelet/app/options/options.go#L219
	// https://github.com/kubernetes/kubernetes/blob/c7f9dd0bafe2f3328148a85eed9c18322b9f308e/cmd/kubelet/app/options/options.go#L397
	SeccompProfileRoot string
}

// Configurator is an InjectionConfigurator that creates SpecPodInjectionConfigurators
// for injecting pods.
type Configurator struct {
	// config is the set of parameters for injecting pods on this host.
	config *Config
}

var _ cri.InjectionConfigurator = &Configurator{}

// NewConfigurator creates a new Configurator.
func NewConfigurator(config *Config) cri.InjectionConfigurator {
	return &Configurator{config: config}
}

func (c *Configurator) ConfigurePodInjection(ctx context.Context,
	logFields []zap.Field, podSandboxConfig *api.PodSandboxConfig) (
	cri.PodInjectionConfigurator, error) {
	return &PodConfigurator{
		podSandboxConfig: podSandboxConfig,
		spec:             &TestSpec, // TODO(rlenglet): Dynamically generate the spec from podSandboxConfig, etc.
		config:           c.config,
	}, nil
}

// PodConfigurator is a PodInjectionConfigurator that translates Istio
// an SidecarInjectionSpec to inject a pod.
type PodConfigurator struct {
	// podSandboxConfig is the CRI configuration for the pod to inject.
	podSandboxConfig *api.PodSandboxConfig

	// spec is the Istio sidecar injection spec for injecting the pod.
	spec *inject.SidecarInjectionSpec

	// config is the set of parameters for injecting pods on this host.
	config *Config
}

var _ cri.PodInjectionConfigurator = &PodConfigurator{}

func (c *PodConfigurator) DNSConfig() *api.DNSConfig {
	dnsConfig := c.spec.DNSConfig
	if dnsConfig == nil {
		return nil
	}
	return &api.DNSConfig{
		Searches: dnsConfig.Searches,
		// ValidatePodDNSConfig ensures that there are no options or nameservers.
	}
}

func (c *PodConfigurator) PortMappings() []*api.PortMapping {
	proxyContainer := c.spec.Containers[0]
	l := len(proxyContainer.Ports)
	if l == 0 {
		return nil
	}

	portMappings := make([]*api.PortMapping, 0, l)
	for _, port := range proxyContainer.Ports {
		protocol := api.Protocol_TCP
		switch port.Protocol {
		case k8s.ProtocolUDP:
			protocol = api.Protocol_UDP
		case k8s.ProtocolSCTP:
			protocol = api.Protocol_SCTP
		}
		portMappings = append(portMappings, &api.PortMapping{
			Protocol:      protocol,
			ContainerPort: port.ContainerPort,
			HostPort:      port.HostPort,
			HostIp:        port.HostIP,
		})
	}
	return portMappings
}

func (c *PodConfigurator) ConfigurePodInjectionForStatus(logFields []zap.Field,
	podSandboxStatus *api.PodSandboxStatus) (*cri.PodInjectionConfig, error) {
	// TODO(rlenglet): Redesign PodInjectionConfig???
	return &cri.PodInjectionConfig{
		IstioInitContainer: *c.TranslateContainer(
			&c.spec.InitContainers[0], podSandboxStatus),
		IstioProxyInjectionEnabled: true, // TODO(rlenglet)
		IstioProxyContainer: *c.TranslateContainer(
			&c.spec.Containers[0], podSandboxStatus),
	}, nil
}

// TranslateContainer translates the given container spec into an InjectedContainerConfig.
func (c *PodConfigurator) TranslateContainer(container *k8s.Container,
	podSandboxStatus *api.PodSandboxStatus) *cri.InjectedContainerConfig {
	// TODO(rlenglet): Translate mounts.
	fixmeTestMounts := []*api.Mount{
		{
			ContainerPath: "/etc/istio/proxy",
			HostPath:      "/tmp/test/etc-istio-proxy", // TODO(rlenglet): Generate that path dynamically.
			// TODO(rlenglet): Set chmod rights to make sure it's writable by user 1337.
		},
		{
			ContainerPath: "/etc/certs/",
			HostPath:      "/tmp/test/etc-certs", // TODO(rlenglet): Generate that path dynamically.
		},
		{
			ContainerPath: "/etc/hosts",
			HostPath:      "/tmp/test/etc-hosts", // TODO(rlenglet): Generate that path and file contents dynamically.
			// TODO(rlenglet): Cf. kubelet's managedHostsFileContent() function.
		},
	}

	return &cri.InjectedContainerConfig{
		Name:       container.Name,
		Image:      container.Image,
		Command:    container.Command,
		Args:       container.Args,
		WorkingDir: container.WorkingDir,
		Envs:       c.TranslateEnv(container.Env, podSandboxStatus),
		Mounts:     fixmeTestMounts,
		Labels: map[string]string{
			"io.kubernetes.container.name": container.Name,
			"io.kubernetes.pod.name":       c.podSandboxConfig.GetMetadata().GetName(),
			"io.kubernetes.pod.namespace":  c.podSandboxConfig.GetMetadata().GetNamespace(),
		},
		Annotations: nil,
		Linux: &api.LinuxContainerConfig{
			Resources: c.TranslateResourceRequirements(&container.Resources),
			SecurityContext: c.TranslateSecurityContext(container.Name,
				container.SecurityContext),
		},
	}
}

// TranslateEnv translates the given environment variables into CRI container environment
// variables.
func (c *PodConfigurator) TranslateEnv(envVars []k8s.EnvVar,
	podSandboxStatus *api.PodSandboxStatus) []*api.KeyValue {
	l := len(envVars)
	if l == 0 {
		return nil
	}

	res := make([]*api.KeyValue, 0, l)
	for _, env := range envVars {
		key := env.Name
		value := env.Value
		if value == "" && env.ValueFrom != nil {
			switch env.ValueFrom.FieldRef.FieldPath {
			case "metadata.name":
				value = c.podSandboxConfig.GetMetadata().GetName()
			case "metadata.namespace":
				value = c.podSandboxConfig.GetMetadata().GetNamespace()
			case "status.hostIP":
				value = c.config.HostIP
			case "status.podIP":
				value = podSandboxStatus.GetNetwork().GetIp()
			}
		}
		res = append(res, &api.KeyValue{Key: key, Value: value})
	}
	return res
}

// TranslateResourceRequirements translates the given resource requirements into CRI container
// resources.
func (c *PodConfigurator) TranslateResourceRequirements(
	resources *k8s.ResourceRequirements) *api.LinuxContainerResources {
	if resources == nil {
		return nil
	}

	// Translate the resource requests and limits similarly to
	// https://github.com/kubernetes/kubernetes/blob/f49fe2a7504c40684573f3e4928e8b0bfdfb2b30/pkg/kubelet/cm/helpers_linux.go

	cpuRequest := resources.Requests.Cpu()
	var cpuShares int64
	if !cpuRequest.IsZero() {
		cpuShares = (cpuRequest.MilliValue() * cpuSharesPerCPU) / milliCPUsPerCPU
	}
	if cpuShares < cpuSharesMin {
		cpuShares = cpuSharesMin
	}

	cpuLimit := resources.Limits.Cpu()
	var cpuPeriod int64
	var cpuQuota int64
	if !cpuLimit.IsZero() {
		cpuPeriod = cpuPeriodDefault
		cpuQuota = (cpuLimit.MilliValue() * cpuPeriod) / milliCPUsPerCPU
		if cpuQuota < cpuQuotaMin {
			cpuQuota = cpuQuotaMin
		}
	}

	memoryLimit := resources.Limits.Memory().Value()

	// Translate the OOM score adjustment similarly to
	// https://github.com/kubernetes/kubernetes/blob/6796645672b4de8d2b6fae780c660926a23e326a/pkg/kubelet/qos/policy.go#L44
	//
	// Since the memory resources reserved for the pod do not account for this container, it is
	// only fair to consider this container as "best effort" re: memory consumption.
	var oomScoreAdj int64 = oomScoreAdjBestEffort

	res := &api.LinuxContainerResources{
		CpuPeriod:          cpuPeriod,
		CpuQuota:           cpuQuota,
		CpuShares:          cpuShares,
		MemoryLimitInBytes: memoryLimit,
		OomScoreAdj:        oomScoreAdj,
		// The following fields are not set by kubelet, so leaving them unset as well:
		// CpusetCpus: ...,
		// CpusetMems: ...,
	}

	return res
}

// TranslateSecurityContext translates the given security context into a CRI container security
// context.
func (c *PodConfigurator) TranslateSecurityContext(
	containerName string, ctx *k8s.SecurityContext) *api.LinuxContainerSecurityContext {
	if ctx == nil {
		return nil
	}

	// Translate the security context similarly to
	// https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/kubelet/kuberuntime/security_context.go#L29

	var capabilities *api.Capability
	if ctx.Capabilities != nil &&
		(len(ctx.Capabilities.Add) > 0 || len(ctx.Capabilities.Drop) > 0) {
		capabilities = &api.Capability{}
		if l := len(ctx.Capabilities.Add); l > 0 {
			capabilities.AddCapabilities = make([]string, l)
			for i, c := range ctx.Capabilities.Add {
				capabilities.AddCapabilities[i] = string(c)
			}
		}
		if l := len(ctx.Capabilities.Drop); l > 0 {
			capabilities.DropCapabilities = make([]string, l)
			for i, c := range ctx.Capabilities.Drop {
				capabilities.DropCapabilities[i] = string(c)
			}
		}
	}

	var seLinuxOptions *api.SELinuxOption
	if ctx.SELinuxOptions != nil {
		seLinuxOptions = &api.SELinuxOption{
			User:  ctx.SELinuxOptions.User,
			Role:  ctx.SELinuxOptions.Role,
			Type:  ctx.SELinuxOptions.Type,
			Level: ctx.SELinuxOptions.Level,
		}
	}

	var runAsUser *api.Int64Value
	if ctx.RunAsUser != nil {
		runAsUser = &api.Int64Value{Value: *ctx.RunAsUser}
	}
	var runAsGroup *api.Int64Value
	if ctx.RunAsGroup != nil {
		runAsGroup = &api.Int64Value{Value: *ctx.RunAsGroup}
	}

	maskedPaths := defaultMaskedPaths
	readonlyPaths := defaultReadonlyPaths
	if ctx.ProcMount != nil && *ctx.ProcMount == k8s.UnmaskedProcMount {
		maskedPaths = nil
		readonlyPaths = nil
	}

	return &api.LinuxContainerSecurityContext{
		Capabilities: capabilities,
		Privileged:   ctx.Privileged != nil && *ctx.Privileged,
		NamespaceOptions: &api.NamespaceOption{
			Network: api.NamespaceMode_POD,
			Pid:     api.NamespaceMode_CONTAINER,
			Ipc:     api.NamespaceMode_POD,
		},
		SelinuxOptions: seLinuxOptions,
		RunAsUser:      runAsUser,
		RunAsGroup:     runAsGroup,
		// The username doesn't need to be set. It's not yet supported in the K8s API, and it
		// seems required only on Windows, cf.
		// https://kubernetes.io/docs/setup/production-environment/windows/intro-windows-in-kubernetes/#security
		RunAsUsername:      "",
		ReadonlyRootfs:     ctx.ReadOnlyRootFilesystem != nil && *ctx.ReadOnlyRootFilesystem,
		SupplementalGroups: nil, // Cannot be set in spec.
		ApparmorProfile:    c.getApparmorProfile(containerName),
		SeccompProfilePath: c.getSeccompProfilePath(containerName),
		NoNewPrivs:         ctx.AllowPrivilegeEscalation != nil && !*ctx.AllowPrivilegeEscalation,
		MaskedPaths:        maskedPaths,
		ReadonlyPaths:      readonlyPaths,
	}
}

// getApparmorProfile returns the Apparmor profile to use for an injected container.
//
// This normally requires users to set one annotation on the pod for each container, cf.
// https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/security/apparmor/helpers.go#L60
// Since pods shouldn't have to refer to injected containers, this is instead determined from the
// annotations added in the injection spec.
func (c *PodConfigurator) getApparmorProfile(containerName string) string {
	return c.spec.PodRedirectAnnot[ApparmorContainerAnnotationKeyPrefix+containerName]
}

// getSeccompProfilePath returns the seccomp profile path to use for an injected container.
//
// This normally requires users to set one annotation on the pod for each container, cf.
// https://github.com/kubernetes/kubernetes/blob/d24fe8a801748953a5c34fd34faa8005c6ad1770/pkg/kubelet/kuberuntime/helpers.go#L209
// Since pods shouldn't have to refer to injected containers, this is instead determined from the
// annotations added in the injection spec.
func (c *PodConfigurator) getSeccompProfilePath(containerName string) string {
	profile, profileOK := c.podSandboxConfig.GetAnnotations()[k8s.SeccompPodAnnotationKey]
	if containerName != "" {
		containerProfile, containerProfileOK := c.spec.PodRedirectAnnot[
			k8s.SeccompContainerAnnotationKeyPrefix+containerName]
		if containerProfileOK {
			profile = containerProfile
			profileOK = containerProfileOK
		}
	}

	if !profileOK {
		return ""
	}

	if strings.HasPrefix(profile, "localhost/") {
		name := strings.TrimPrefix(profile, "localhost/")
		fname := filepath.Join(c.config.SeccompProfileRoot, filepath.FromSlash(name))
		return "localhost/" + fname
	}

	return profile
}
