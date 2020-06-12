package utils

import (
	"os"
)

func GetKubernetesProvider() string {
	provider, found := os.LookupEnv(KubernetesProviderEnvVar)
	if !found {
		return ""
	}
	return provider
}

// IsOpenShiftCluster returns true if this is a openshift cluster
func IsOpenShiftCluster() bool {
	if GetKubernetesProvider() == OpenShiftProvider {
		return true
	}
	return false
}

func GetAdminPassword() string {
	password, found := os.LookupEnv(AdminPasswordEnvVar)
	if !found {
		return DefaultAdminPassword
	}
	return password
}

func GetKubernetesApiServer() string {
	server, found := os.LookupEnv(KubernetesServiceHostEnvVar)
	if !found {
		return ""
	}
	return server
}

func GetKubernetesApiPort() string {
	port, found := os.LookupEnv(KubernetesServicePortEnvVar)
	if !found {
		return "6443"
	}
	return port
}
