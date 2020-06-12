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

	// KubernetesServiceHostEnvVar variable defines the env variable
	// KUBERNETES_SERVICE_HOST, which is used to indicate k8s api server
	// host address
	KubernetesServiceHostEnvVar = "KUBERNETES_SERVICE_HOST"

	// KubernetesServicePortEnvVar variable defines the env variable
	// KUBERNETES_SERVICE_PORT, which is used to indicate k8s api server
	// port address
	KubernetesServicePortEnvVar = "KUBERNETES_SERVICE_PORT"

	// Value for k8s provider as OpenShift
	OpenShiftProvider = "OpenShift"

	// Default value used for password
	DefaultAdminPassword = "atsgen"
)
