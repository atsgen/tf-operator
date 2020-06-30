package values

const (
	// TFDefaultDeployment - Tungsten fabric default deployment name
	TFDefaultDeployment    = "default"

	// TFNamespace - Tungsten Fabric namespace
	TFNamespace            = "tungsten"

	// TFReleaseTag - Default release tag to be used
	// ensure operator uses the latest supported images for this
	// deployment, ensure upgrade by just upgrading the operator
	// version
	TFReleaseTag           = "auto"

	// TFCurrentRelease - current tungsten fabric release this
	// operator will deploy, if the release tag is TFReleaseTag
	// or not defined
	TFCurrentRelease       = "R2003-latest"

	// OpenShiftNetworkConfig - OpenShift network config name
	OpenShiftNetworkConfig     = "cluster"

	// OpenShiftAtsgenCni - OpenShift Atsgen CNI label
	OpenShiftAtsgenCni         = "atsgenTungsten"

	// DefaultCniBinDir - default cni bin directory for K8s
	DefaultCniBinDir           = "/opt/cni/bin"

	// DefaultCniConfDir - default cni conf directory for K8s
	DefaultCniConfDir          = "/etc/cni"

	// OpenShiftCniBinDir - cni bin directory for OpenShift
	OpenShiftCniBinDir         = "/var/lib/cni/bin"

	// OpenShiftCniConfDir - cni conf directory for OpenShift
	OpenShiftCniConfDir        = "/etc/kubernetes/cni/"

	// OpenShiftMultusConfDir - cni conf directory for OpenShift,
	// when multus is enabled
	OpenShiftMultusConfDir     = "/var/run/multus/cni/"
)
