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
