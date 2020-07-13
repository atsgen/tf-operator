package utils

import (
	"fmt"
	"os"

	"github.com/atsgen/tf-operator/pkg/values"
)

// GetKubernetesProvider returns the K8s provider for this deployment
func GetKubernetesProvider() string {
	provider, found := os.LookupEnv(KubernetesProviderEnvVar)
	if !found {
		return ""
	}
	return provider
}

// IsOpenShiftCluster returns true if this is a openshift cluster
func IsOpenShiftCluster() bool {
	return (GetKubernetesProvider() == OpenShiftProvider)
}

// GetAdminPassword returns the admin password configured
func GetAdminPassword() string {
	password, found := os.LookupEnv(AdminPasswordEnvVar)
	if !found {
		return DefaultAdminPassword
	}
	return password
}

// SetAdminPassword sets the admin password for later use
func SetAdminPassword(password string) {
	_ = os.Setenv(AdminPasswordEnvVar, password)
}

// GetKubernetesAPIServer returns out of cluster configuration for K8S API
func GetKubernetesAPIServer() string {
	server, found := os.LookupEnv(KubernetesServiceHostEnvVar)
	if !found {
		return ""
	}
	return server
}

// GetKubernetesAPIPort returns K8S API port
func GetKubernetesAPIPort() string {
	port, found := os.LookupEnv(KubernetesServicePortEnvVar)
	if !found {
		return "6443"
	}
	return port
}

// IsOpenShiftMultusEnabled returns if multus is enabled
func IsOpenShiftMultusEnabled() (bool, error) {
	status, found := os.LookupEnv(OpenShiftMultusStatusEnvVar)
	if !found {
		return false, fmt.Errorf("does not exist")
	}

	if status != "enabled" {
		return false, nil
	}
	return true, nil
}

// SetOpenShiftMultusStatus sets multus status as env variable
// for consumption later
func SetOpenShiftMultusStatus(enabled bool) {
	if enabled {
		_ = os.Setenv(OpenShiftMultusStatusEnvVar, "enabled")
	} else {
		_ = os.Setenv(OpenShiftMultusStatusEnvVar, "disabled")
	}
}

// GetContainerRegistry returns the container registry for this
// deployment
func GetContainerRegistry() string {
	registry, found := os.LookupEnv(ContainerRegistryEnvVar)
	if !found {
		return DefaultContainerRegistry
	}

	return registry
}

// GetContainerPrefix returns the container prefix to be used
// for this deployment
func GetContainerPrefix() string {
	prefix, found := os.LookupEnv(ContainerPrefixEnvVar)
	if !found {
		if IsOpenShiftCluster() {
			// unless overriden for openshift cluster
			// we cannot use contrail prefix images
			return ContainerPrefixTungsten
		}
		return ContainerPrefixContrail
	}
	return prefix
}

// GetReleaseTag returns the operational release tag give
// the configured release tag
func GetReleaseTag(configTag string) string {
	switch configTag {
	case "":
		fallthrough
	case values.TFReleaseTag:
		return values.TFCurrentRelease
	}
	return configTag
}

// IsTungstenFabricHADisabled returns the status for Tungsten
// fabric HA configuration
func IsTungstenFabricHADisabled() bool {
	_, found := os.LookupEnv(DisableTungstenHAEnvVar)
	if !found {
		// currently we support HA only for OpenShift cluster
		return !IsOpenShiftCluster()
	}
	// we don't care the value defined if the environment variable
	// exists it referes to disabling HA
	return true
}

// GetOperatorNamespace returns the namespace in which operator is
// running
func GetOperatorNamespace() string {
	namespace, found := os.LookupEnv(OperatorNamespaceEnvVar)
	if !found {
		return ""
	}
	return namespace
}

// IsResourceHackDisabled returns the status for resource hack for limiting
// the ram usage
func IsResourceHackDisabled() string {
	_, found := os.LookupEnv(DisableResourceHackEnvVar)
	if !found {
		return FalseStr
	}

	// we don't care the value, if the env variable is define
	// it referes to disabling the resource limit hack
	return TrueStr
}
