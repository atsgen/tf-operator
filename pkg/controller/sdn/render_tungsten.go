package sdn

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/types"
	uns "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/atsgen/tf-operator/pkg/apply"
	"github.com/atsgen/tf-operator/pkg/render"
	"github.com/atsgen/tf-operator/pkg/utils"
	"github.com/atsgen/tf-operator/pkg/values"
)

func updateCNIInfo(data *render.RenderData, config *tungstenv1alpha1.CNIConfigType) {
	if config.ClusterName == "" {
		data.Data["KUBERNETES_CLUSTER_NAME"] = "k8s"
	} else {
		data.Data["KUBERNETES_CLUSTER_NAME"] = config.ClusterName
	}

	if config.UseHostNewtorkService {
		data.Data["KUBERNETES_HOST_NETWORK_SERVICE"] = "true"
	} else {
		data.Data["KUBERNETES_HOST_NETWORK_SERVICE"] = "false"
	}

	if utils.IsOpenShiftCluster() {
		data.Data["CNI_BIN_DIR"] = values.OpenShiftCniBinDir
		multusEnabled, _ := utils.IsOpenShiftMultusEnabled()
		if multusEnabled {
			data.Data["CNI_CONF_DIR"] = values.OpenShiftMultusConfDir
		} else {
			data.Data["CNI_CONF_DIR"] = values.OpenShiftCniConfDir
		}
	} else {
		data.Data["CNI_BIN_DIR"] = values.DefaultCniBinDir
		data.Data["CNI_CONF_DIR"] = values.DefaultCniConfDir
	}

	data.Data["KUBERNETES_POD_SUBNETS"] = config.PodNetwork.Cidr
	data.Data["KUBERNETES_SERVICE_SUBNETS"] = config.ServiceNetwork.Cidr
	data.Data["KUBERNETES_IP_FABRIC_SUBNETS"] = config.IpFabricNetwork.Cidr

	switch config.IpForwarding {
	case IPForwardingEnabled:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "true"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "false"
	case IPForwardingSnat:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "false"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "true"
	default:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "false"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "false"
	}
}

func updateDatapathInfo(data *render.RenderData) {
	// TODO(prabhjot) some bug with encryption, disable for now
	data.Data["VROUTER_ENCRYPTION"] = "FALSE"
	data.Data["VROUTER_GATEWAY"] = ""
	data.Data["DPDK_UIO_DRIVER"] = "igb_uio"

	if utils.IsOpenShiftCluster() {
		// we don't support building KMOD for openshift
		data.Data["TUNGSTEN_KMOD"] = "init"
	} else {
		data.Data["TUNGSTEN_KMOD"] = "build"
	}
}

func updateContainerInfo(data *render.RenderData, cr *tungstenv1alpha1.SDN) {
	data.Data["CONTAINER_REGISTRY"] = utils.GetContainerRegistry()
	data.Data["CONTAINER_PREFIX"] = utils.GetContainerPrefix()
	data.Data["CONTAINER_TAG"] = utils.GetReleaseTag(cr.Spec.ReleaseTag)
}

func solicitData(data *render.RenderData, cr *tungstenv1alpha1.SDN, nodes *NodeList) {
	var controllerNodes string
	for ip, _ := range controllerIPs {
		if controllerNodes == "" {
			controllerNodes = ip
		} else {
			controllerNodes = controllerNodes + "," + ip
		}
	}

	data.Data["K8S_PROVIDER"] = utils.GetKubernetesProvider()
	data.Data["TF_NAMESPACE"] = values.TFNamespace
	data.Data["AAA_MODE"] = "no-auth"
	data.Data["ADMIN_PASSWORD"] = utils.GetAdminPassword()
	data.Data["ANALYTICS_ALARM_NODES"] = controllerNodes
	data.Data["ANALYTICS_API_VIP"] = ""
	data.Data["ANALYTICSDB_NODES"] = controllerNodes
	data.Data["ANALYTICS_NODES"] = controllerNodes
	data.Data["ANALYTICS_SNMP_NODES"] = controllerNodes
	data.Data["AUTH_MODE"] = "noauth"
	data.Data["CLOUD_ORCHESTRATOR"] = "kubernetes"
	data.Data["CONFIG_API_VIP"] = ""
	data.Data["CONFIGDB_NODES"] = controllerNodes
	data.Data["CONFIG_NODES"] = controllerNodes

	// update container information
	updateContainerInfo(data, cr)

	updateCNIInfo(data, &cr.Spec.CNIConfig)

	updateDatapathInfo(data)

	data.Data["CONTROLLER_NODES"] = controllerNodes
	data.Data["CONTROL_NODES"] = controllerNodes
	data.Data["JVM_EXTRA_OPTS"] = "-Xms1g -Xmx2g"
	data.Data["KAFKA_NODES"] = controllerNodes
	data.Data["KUBERNETES_API_SECURE_PORT"] = utils.GetKubernetesAPIPort()
	apiServer := utils.GetKubernetesAPIServer()
	if apiServer == "" {
		apiServer = nodes.DefultApiServer
	}
	data.Data["KUBERNETES_API_SERVER"] = apiServer
	data.Data["KUBERNETES_PUBLIC_FIP_POOL"] = ""
	data.Data["TUNGSTEN_IMAGE_PULL_SECRET"] = ""
	data.Data["LOG_LEVEL"] = "SYS_NOTICE"
	data.Data["METADATA_PROXY_SECRET"] = "tungsten"
	data.Data["PHYSICAL_INTERFACE"] = ""
	data.Data["RABBITMQ_NODE_PORT"] = "5673"
	data.Data["RABBITMQ_NODES"] = controllerNodes
	data.Data["WEBUI_NODES"] = controllerNodes
	data.Data["WEBUI_VIP"] = ""
	data.Data["ZOOKEEPER_PORT"] = "2181"
	data.Data["ZOOKEEPER_PORTS"] = "2888:3888"
}

func checkAndRenderStage(r *ReconcileSDN, cr *tungstenv1alpha1.SDN, data *render.RenderData, stage string) (string, error) {
	objs := []*uns.Unstructured{}

	manifests, err := render.RenderDir(filepath.Join("/bindata", "tungsten/", stage), data)
	if err != nil {
		log.Info("Failed to render yaml files " + err.Error());
		return "", err
	}

	objs = append(objs, manifests...)
	if waitForRollout(stage) {
		rolloutPending := false
		var daemonSets []*appsv1.DaemonSet = nil
		for _, obj := range objs {
			if obj.GetAPIVersion() == "apps/v1" && obj.GetKind() == "DaemonSet" {
				ds := &appsv1.DaemonSet{}
				err := r.client.Get(context.TODO(),
					types.NamespacedName{Namespace: obj.GetNamespace(),
							Name: obj.GetName()},
					ds)
				if err != nil {
					// we still haven't created this resource
					rolloutPending = true
					daemonSets = nil
					break
				}
				daemonSets = append(daemonSets, ds)
			}
		}
		if !rolloutPending {
			for _, ds := range daemonSets {
				if ds.Status.DesiredNumberScheduled != ds.Status.NumberAvailable {
					// we have not reached the required number of nodes
					return TFOperatorObjectUpdating, nil
				}
			}
			// all the daemon sets for this stage are up and running
			return TFOperatorObjectDeployed, nil
		}
	}

	if utils.IsOpenShiftCluster() && (stage == "base") {
                // cluster is running for openshift, load objects needed for openshift
		manifests, err := render.RenderDir(filepath.Join("/bindata", "openshift/"), data)
		if err != nil {
			log.Info("Failed to render yaml files " + err.Error());
			return "", err
		}
		objs = append(objs, manifests...)
	}

	for _, obj := range objs {
		if err := controllerutil.SetControllerReference(cr, obj, r.scheme); err!= nil {
			log.Info(err.Error())
			return "", err
		}
		if err := apply.ApplyObject(context.TODO(), r.client, obj); err != nil {
			log.Info(err.Error())
			err = errors.Wrapf(err, "could not apply (%s) %s/%s", obj.GroupVersionKind(), obj.GetNamespace(), obj.GetName())
			return "", err
		}
	}

	if waitForRollout(stage) {
		// let elements roll-out and we will reconcile again to proceed
		// to the next stage
		return TFOperatorObjectUpdating, nil
	}
	return TFOperatorObjectDeployed, nil
}

func renderTungstenFabric(r *ReconcileSDN, cr *tungstenv1alpha1.SDN, nodes *NodeList) (string, error) {
	data := render.MakeRenderData()
	solicitData(&data, cr, nodes)

	stage := cr.Status.Stage
	if stage == "" {
		stage = getNextStage(stage)
	}
	for stage != "deployed" {
		status, err := checkAndRenderStage(r, cr, &data, stage)
		if status == TFOperatorObjectDeployed {
			// go to next stage
			stage = getNextStage(stage)
			r.updateStage(cr, stage)
		} else {
			return status, err
		}
	}
	return TFOperatorObjectDeployed, nil
}
