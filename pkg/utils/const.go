package utils

const (
	// KubernetesProvider variable defines the env variable K8S_PROVIDER
	// which is currently used to indicate operator, if openshift constructs
	// needs to be enabled
	KubernetesProviderEnvVar = "K8S_PROVIDER"

	// AdminPasswordEnvVar variable defines the env variable ADMIN_PASSWORD
	// which holds the value of secret to be used as Tugnsten Fabric Admin
	// password
	AdminPasswordEnvVar = "ADMIN_PASSWORD"

	// Value for k8s provider as OpenShift
	OpenShiftProvider = "OpenShift"

	// Default value used for password
	DefaultAdminPassword = "atsgen"
)
