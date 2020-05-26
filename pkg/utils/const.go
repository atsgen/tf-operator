package utils

const (
	// KubernetesProvider variable defines the env variable K8S_PROVIDER
	// which is currently used to indicate operator, if openshift constructs
	// needs to be enabled
	KubernetesProviderEnvVar = "K8S_PROVIDER"

	// Value for k8s provider as OpenShift
	OpenShiftProvider = "OpenShift"
)
