package sdn

import (
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	"github.com/atsgen/tf-operator/pkg/values"
)

func Validate(cr *tungstenv1alpha1.SDN) (state string, description string) {
	if cr.Name != values.TFDefaultDeployment {
		return TFOperatorObjectIgnored, ("Tungsten CNI other than name: " + values.TFDefaultDeployment + ", are not processed")
	}

	s := TFOperatorNotSupported
	d := "un-supported release tag, deployment will proceed as un-supported"
	if cr.Spec.ReleaseTag == values.TFReleaseTag {
		s = TFOperatorObjectDeployed
		d = ""
	}
	return s, d
}
