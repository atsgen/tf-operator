package tungstencni

import (
	"github.com/atsgen/tf-operator/pkg/utils"
)

func getNextStage(stage string) string {
	switch stage {
	case "":
		return "base"
	case "base":
		return "controller"
	case "controller":
		return "datapath"
	case "datapath":
		return "analytics"
	case "analytics":
		return "webui"
	}
	return "deployed"
}

func waitForRollout(stage string) bool {
	if !utils.IsOpenShiftCluster() {
		// we will need to wait for rollouts only in an HA based
		// deployments
		return false
	}

	switch stage {
	case "base", "controller":
		return true
	}
	return false
}
