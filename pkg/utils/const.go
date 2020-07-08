package utils

const (
	// KubernetesProviderEnvVar variable defines the env variable K8S_PROVIDER
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

	// OpenShiftMultusStatusEnvVar variable defines the env variable
	// OPENSHIFT_MULTUS, which is used to indicate the whether multus
	// is enabled for OpenShift or not
	OpenShiftMultusStatusEnvVar = "OPENSHIFT_MULTUS"

	// ContainerRegistryEnvVar variable defines the env variable
	// CONTAINER_REGISTRY, which is used to indicate container registry
	// to use
	ContainerRegistryEnvVar = "CONTAINER_REGISTRY"

	// ContainerPrefixEnvVar variable defines the env variable
	// CONTAINER_PREFIX, which is used to indicate container prefix
	// to use, this is needed to toggle between tungsten and contrail
	// will be removed once contrail references are completely removed
	ContainerPrefixEnvVar = "CONTAINER_PREFIX"

	// DisableTungstenHAEnvVar variable defines the env variable
	// DISABLE_TUNGSTEN_HA, which is used to indicate disabling of
	// HA for OpenShift deployment
	DisableTungstenHAEnvVar = "DISABLE_TUNGSTEN_HA"

	// OpenShiftProvider - Value for k8s provider as OpenShift
	OpenShiftProvider = "OpenShift"

	// DefaultAdminPassword - Default value used for password
	DefaultAdminPassword = "atsgen"

	// DefaultContainerRegistry - default value for container registry
	DefaultContainerRegistry = "atsgen"

	// ContainerPrefixContrail - images to use contrail as prefix
	ContainerPrefixContrail = "contrail"

	// ContainerPrefixTungsten - images to use tungsten as prefix
	ContainerPrefixTungsten = "tungsten"
)
