package controller

import (
	"github.com/atsgen/tf-operator/pkg/controller/tungstencni"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, tungstencni.Add)
}
