package apis

import (
	"github.com/openshift/api/config/v1"
	"github.com/atsgen/tf-operator/pkg/utils"
)

func init() {
	// register only for openshift cluster
	if utils.IsOpenShiftCluster() {
		// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
		AddToSchemes = append(AddToSchemes, v1.AddToScheme)
	}
}
