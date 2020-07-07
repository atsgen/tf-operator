package sdn

import (
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	"github.com/atsgen/tf-operator/pkg/utils"
	"github.com/atsgen/tf-operator/pkg/values"
)

func Validate(cr *tungstenv1alpha1.SDN) (state string, description string) {
	if cr.Name != values.TFDefaultDeployment {
		return TFOperatorObjectIgnored, ("Tungsten CNI other than name: " + values.TFDefaultDeployment + ", are not processed")
	}

	// Poor man approach to wait for OpenShift Network config to merge
	// into SDN CR
	if utils.IsOpenShiftCluster() && cr.Spec.CNIConfig.ClusterName == "" {
		return TFOperatorObjectIgnored, "Waiting for OpenShift network config"
	}

	s := TFOperatorNotSupported
	d := "un-supported release tag, deployment will proceed as un-supported"

	// TODO(Prabhjot) will need a separate util function to
	// validate release tag
	if utils.GetReleaseTag(cr.Spec.ReleaseTag) == values.TFCurrentRelease {
		s = TFOperatorObjectDeployed
		d = ""
	}
	return s, d
}
