package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

// Define PodNetwork parameters for Tugnsten Fabric
type PodNetworkType struct {
	// pod network CIDR
	Cidr string `json:"cidr,omitempty"`
}

// Define ServiceNetwork parameters for Tugnsten Fabric
type ServiceNetworkType struct {
	// service network CIDR
	Cidr string `json:"cidr,omitempty"`
}

// Define IpFabricNetwork parameters for Tugnsten Fabric
type IpFabricNetworkType struct {
	// ip fabric network CIDR
	Cidr string `json:"cidr,omitempty"`
}

// +kubebuilder:validation:Enum=snat;enable
type IpForwardingType string

// CNIConfigType placeholder for Tungsten Fabric CNI configuration
type CNIConfigType struct {
	// cluster name
	ClusterName string  `json:"clusterName,omitempty"`

	// pod network parameters
	PodNetwork PodNetworkType `json:"podNetwork,omitempty"`

	// service network parameters
	ServiceNetwork ServiceNetworkType `json:"serviceNetwork,omitempty"`

	// ip fabric network parameters
	IpFabricNetwork IpFabricNetworkType `json:"ipFabricNetwork,omitempty"`

	// ip fabric forwarding, supported value enable, snat,
	// where as empty field means ip fabric forwarding disabled
	IpForwarding IpForwardingType  `json:"ipForwarding,omitempty"`

	// use host network services
	UseHostNewtorkService bool  `json:"useHostNewtorkService,omitempty"`
}

// DatapathConfigType placeholder for Tungsten Fabric Datapath related config
type DatapathConfigType struct {
	// UseVrouter as datpath for CNI
	UseVrouter bool   `json:"useVrouter,omitempty"`
}

// Define the desired SDN deployment parameters
type SDNSpec struct {
	// ReleaseTag for the container images used empty release tag would
	// assume automatically move to the latest image tag supported
	ReleaseTag string `json:"releaseTag,omitempty"`

	// AdminPassword password for Tungsten Fabric Controller admin
	AdminPassword string `json:"adminPassword,omitempty"`

	// CNIConfig supplies configuration used for Tungsten Fabric CNI
	CNIConfig CNIConfigType  `json:"cniConfig,omitempty"`

	// DatapathConfig supplies configuration for Tungsten Fabric Datapath
	DatapathConfig DatapathConfigType `json:"datapathConfig,omitempty"`
}

// SDNStatus defines the observed state of SDN
type SDNStatus struct {
	// state of deployment
	State string `json:"state,omitempty"`
	Error string `json:"error,omitempty"`
	ReleaseTag string `json:"releaseTag,omitempty"`
	Stage string `json:"stage,omitempty"`
	Controllers []string `json:"controllers,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SDN is the Schema for the sdns API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sdns,scope=Cluster
// +kubebuilder:printcolumn:name="Release",type=string,JSONPath=`.status.releaseTag`
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SDN struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SDNSpec   `json:"spec,omitempty"`
	Status SDNStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SDNList contains a list of SDN
type SDNList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SDN `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SDN{}, &SDNList{})
}
