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

// Define the desired TungstenCNI deployment parameters
type TungstenCNISpec struct {
	// release tag for the container images used
	ReleaseTag string `json:"releaseTag"`
	// use vrouter as datpath for CNI
	UseVrouter bool   `json:"useVrouter,omitempty"`
	// pod network parameters
	PodNetwork PodNetworkType `json:"podNetwork,omitempty"`
	// service network parameters
	ServiceNetwork ServiceNetworkType `json:"serviceNetwork,omitempty"`
}

// TungstenCNIStatus defines the observed state of TungstenCNI
type TungstenCNIStatus struct {
	// state of deployment
	State string `json:"state,omitempty"`
	Error string `json:"error,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TungstenCNI is the Schema for the tungstencnis API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=tungstencnis,scope=Cluster
// +kubebuilder:printcolumn:name="Release",type=string,JSONPath=`.spec.releasetag`
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type TungstenCNI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TungstenCNISpec   `json:"spec,omitempty"`
	Status TungstenCNIStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TungstenCNIList contains a list of TungstenCNI
type TungstenCNIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TungstenCNI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TungstenCNI{}, &TungstenCNIList{})
}
