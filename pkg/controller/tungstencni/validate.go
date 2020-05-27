package tungstencni

import (
	"strings"

	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	"github.com/atsgen/tf-operator/pkg/values"
)

func Validate(cr *tungstenv1alpha1.TungstenCNI) (state string, description string) {
	if cr.Name != values.TF_OPERATOR_CONFIG {
		return TF_OPERATOR_OBJECT_IGNORED, ("Tungsten CNI other than name: " + values.TF_OPERATOR_CONFIG + ", are not processed")
	}

	s := TF_OPERATOR_OBJECT_NOT_SUPPORTED
	d := "un-supported release tag, deployment will proceed as un-supported"
	if strings.Contains(cr.Spec.ReleaseTag, "R2003") {
		s = TF_OPERATOR_OBJECT_DEPLOYED
		d = ""
	}
	return s, d
}
