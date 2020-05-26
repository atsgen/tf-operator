package controller

import (
	"github.com/atsgen/tf-operator/pkg/controller/network"
	"github.com/atsgen/tf-operator/pkg/utils"
)

func init() {
	// register only for openshift cluster
	if utils.IsOpenShiftCluster() {
		// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
		AddToManagerFuncs = append(AddToManagerFuncs, network.Add)
	}
}
