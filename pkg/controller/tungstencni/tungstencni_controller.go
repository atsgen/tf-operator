package tungstencni

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	uns "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/atsgen/tf-operator/pkg/apply"
	"github.com/atsgen/tf-operator/pkg/render"
	"github.com/atsgen/tf-operator/pkg/utils"
	"github.com/atsgen/tf-operator/pkg/values"
)

var log = logf.Log.WithName("controller_tungstencni")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// TODO(Prabhjot) need to fix taking parameters from CNI Object
func (r *ReconcileTungstenCNI) renderTungstenFabricCNI(cr *tungstenv1alpha1.TungstenCNI) error {
	objs := []*uns.Unstructured{}

	nodes, e := FetchNodeList(r.client)

	if e != nil {
		return e
	}

	computeRole := NODE_ROLE_VPP
	if cr.Spec.UseVrouter {
		computeRole = NODE_ROLE_VROUTER
	}
	agentLabels := []string{computeRole}
	allLabels := []string{computeRole,
				NODE_ROLE_ANALYTICS,
				NODE_ROLE_ANALYTICS_ALARM,
				NODE_ROLE_ANALYTICS_SNMP,
				NODE_ROLE_ANALYTICS_DB,
				NODE_ROLE_CONFIG,
				NODE_ROLE_CONFIG_DB,
				NODE_ROLE_CONTROL,
				NODE_ROLE_WEBUI}
	for _, name := range nodes.WorkerNodes {
		// enable agent for all nodes
		e = SetNodeLabels(r.client, name, agentLabels)
		if e != nil {
			return e
		}
	}

	for _, name := range nodes.MasterNodes {
		// enable all labels for master nodes
		e = SetNodeLabels(r.client, name, allLabels)
		if e != nil {
			return e
		}
	}

	data := render.MakeRenderData()
	data.Data["K8S_PROVIDER"] = utils.GetKubernetesProvider()
	data.Data["TF_NAMESPACE"] = values.TF_NAMESPACE
	data.Data["AAA_MODE"] = "no-auth"
	data.Data["ADMIN_PASSWORD"] = utils.GetAdminPassword()
	data.Data["ANALYTICS_ALARM_NODES"] = nodes.MasterNodesStr
	data.Data["ANALYTICS_API_VIP"] = ""
	data.Data["ANALYTICSDB_NODES"] = nodes.MasterNodesStr
	data.Data["ANALYTICS_NODES"] = nodes.MasterNodesStr
	data.Data["ANALYTICS_SNMP_NODES"] = nodes.MasterNodesStr
	data.Data["AUTH_MODE"] = "noauth"
	data.Data["CLOUD_ORCHESTRATOR"] = "kubernetes"
	data.Data["CONFIG_API_VIP"] = ""
	data.Data["CONFIGDB_NODES"] = nodes.MasterNodesStr
	data.Data["CONFIG_NODES"] = nodes.MasterNodesStr
	data.Data["CONTAINER_REGISTRY"] = "atsgen"
	data.Data["CONTAINER_TAG"] = cr.Spec.ReleaseTag
	if utils.IsOpenShiftCluster() {
		// we don't support building KMOD for openshift
		data.Data["TUNGSTEN_KMOD"] = "init"
	} else {
		data.Data["TUNGSTEN_KMOD"] = "build"
	}
	data.Data["CONTROLLER_NODES"] = nodes.MasterNodesStr
	data.Data["CONTROL_NODES"] = nodes.MasterNodesStr
	data.Data["JVM_EXTRA_OPTS"] = "-Xms1g -Xmx2g"
	data.Data["KAFKA_NODES"] = nodes.MasterNodesStr
	data.Data["KUBERNETES_API_SECURE_PORT"] = "6443"
	data.Data["KUBERNETES_API_SERVER"] = nodes.MasterNodesStr
	data.Data["KUBERNETES_PUBLIC_FIP_POOL"] = ""
	data.Data["TUNGSTEN_IMAGE_PULL_SECRET"] = ""
	data.Data["LOG_LEVEL"] = "SYS_NOTICE"
	data.Data["METADATA_PROXY_SECRET"] = "tungsten"
	data.Data["PHYSICAL_INTERFACE"] = ""
	data.Data["RABBITMQ_NODE_PORT"] = "5673"
	data.Data["RABBITMQ_NODES"] = nodes.MasterNodesStr
	data.Data["VROUTER_GATEWAY"] = ""
	data.Data["WEBUI_NODES"] = nodes.MasterNodesStr
	data.Data["WEBUI_VIP"] = ""
	data.Data["ZOOKEEPER_PORT"] = "2181"
	data.Data["ZOOKEEPER_PORTS"] = "2888:3888"
	data.Data["DPDK_UIO_DRIVER"] = "igb_uio"
	data.Data["KUBERNETES_POD_SUBNETS"] = cr.Spec.PodNetwork.Cidr
	data.Data["KUBERNETES_SERVICE_SUBNETS"] = cr.Spec.ServiceNetwork.Cidr

	manifests, err := render.RenderDir(filepath.Join("/bindata", "tungsten/"), &data)
	if err != nil {
		log.Info("Failed to render yaml files " + err.Error());
		return err
	}

	objs = append(objs, manifests...)
	if utils.IsOpenShiftCluster() {
                // cluster is running for openshift, load objects needed for openshift
		manifests, err := render.RenderDir(filepath.Join("/bindata", "openshift/"), &data)
		if err != nil {
			log.Info("Failed to render yaml files " + err.Error());
			return err
		}
		objs = append(objs, manifests...)
	}

	for _, obj := range objs {
		if err := controllerutil.SetControllerReference(cr, obj, r.scheme); err!= nil {
			log.Info(err.Error())
			return err
		}
		if err := apply.ApplyObject(context.TODO(), r.client, obj); err != nil {
			log.Info(err.Error())
			err = errors.Wrapf(err, "could not apply (%s) %s/%s", obj.GroupVersionKind(), obj.GetNamespace(), obj.GetName())
			return err
		}
	}
	return nil
}

// Add creates a new TungstenCNI Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTungstenCNI{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("tungstencni-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource TungstenCNI
	err = c.Watch(&source.Kind{Type: &tungstenv1alpha1.TungstenCNI{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner TungstenCNI
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tungstenv1alpha1.TungstenCNI{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTungstenCNI implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTungstenCNI{}

// ReconcileTungstenCNI reconciles a TungstenCNI object
type ReconcileTungstenCNI struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TungstenCNI object and makes changes based on the state read
// and what is in the TungstenCNI.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTungstenCNI) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling TungstenCNI")

	// Fetch the TungstenCNI instance
	instance := &tungstenv1alpha1.TungstenCNI{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	s,d := Validate(instance)
	if s == TF_OPERATOR_OBJECT_IGNORED {
		log.Info("Error!!! Ignoring tf-operator " + request.Name)
		r.updateStatus(instance, s, d)
		return reconcile.Result{}, nil
	}

	err = r.renderTungstenFabricCNI(instance)
	if err != nil {
		log.Error(err, "failed to reconcile")
		return reconcile.Result{}, err
	}

	r.updateStatus(instance, s, d)
	log.Info("reconcile completed: Tungsten CNI " + instance.Name + " Updated")
	return reconcile.Result{}, nil
}

func (r *ReconcileTungstenCNI) updateStatus(cr *tungstenv1alpha1.TungstenCNI, state string, msg string) error {
        if (cr.Status.State == state && cr.Status.Error == msg) {
		// No update required
		return nil
	}
	cr.Status.State = state
	cr.Status.Error = msg
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		log.Error(err, "failed to update TungstenCNI status")
		return err
	}
	return nil
}
