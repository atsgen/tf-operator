package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

// Define the desired TungstenCNI deployment parameters
type TungstenCNISpec struct {
	// release tag for the container images used
	ReleaseTag string `json:"releasetag"`
	// use vrouter as datpath for CNI
	UseVrouter bool   `json:"usevrouter,omitempty"`
}

// TungstenCNIStatus defines the observed state of TungstenCNI
type TungstenCNIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TungstenCNI is the Schema for the tungstencnis API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=tungstencnis,scope=Cluster
// +kubebuilder:printcolumn:name="Release",type=string,JSONPath=`.spec.releasetag`
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
