package utils

import (
	"fmt"
	"os"
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

// GetAdminPassword returns the admin password supplied by installation process
func GetAdminPassword() string {
	password, found := os.LookupEnv(AdminPasswordEnvVar)
	if !found {
		return DefaultAdminPassword
	}
	return password
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
