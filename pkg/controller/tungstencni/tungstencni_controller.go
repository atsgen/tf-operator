package tungstencni

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	uns "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/record"

	ocv1 "github.com/openshift/api/operator/v1"

	"github.com/atsgen/tf-operator/pkg/apply"
	"github.com/atsgen/tf-operator/pkg/render"
	"github.com/atsgen/tf-operator/pkg/utils"
	"github.com/atsgen/tf-operator/pkg/values"
)

var log = logf.Log.WithName("controller_tungstencni")

var controllerIPs map[string]bool = make(map[string]bool)

/**
 * USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
 * business logic.  Delete these comments after modifying this file.*
 */

func updateIPForwarding(data *render.RenderData, cr *tungstenv1alpha1.TungstenCNI) {
	switch cr.Spec.IpForwarding {
	case IP_FORWARDING_ENABLED:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "true"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "false"
	case IP_FORWARDING_SNAT:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "false"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "true"
	default:
		data.Data["KUBERNETES_IP_FABRIC_FORWARDING"] = "false"
		data.Data["KUBERNETES_IP_FABRIC_SNAT"] = "false"
	}
}

func (r *ReconcileTungstenCNI) renderTungstenFabricCNI(cr *tungstenv1alpha1.TungstenCNI) (bool, error) {
	objs := []*uns.Unstructured{}

	nodes, e := FetchNodeList(r.client)

	if e != nil {
		return false, e
	}

	// check if we have already identified Controller IPs
	if len(cr.Status.Controllers) == 0 {
		if utils.IsOpenShiftCluster() && len(nodes.MasterNodes) < 3 {
			// return from here we will get notified when a new
			// node is available
			r.recorder.Event(cr, corev1.EventTypeNormal,
				TF_OPERATOR_OBJECT_PENDING,
				fmt.Sprintf("waiting for master node discovery (got %d/%d)", len(nodes.MasterNodes), 3))
			return false, nil
		}
		i := 0
		for ip, _ := range nodes.MasterNodes {
			controllerIPs[ip] = true
			i++
			// we allow only 3 controller nodes for now
			// TODO(prabhjot) will need to consider making this
			// configurable/dynamic
			if i == 3 {
				break
			}
		}
		// commit controller ips to status before going any further
		err := r.updateControllerIPs(cr)
		if err != nil {
			return false, err
		}

		r.recorder.Event(cr, corev1.EventTypeNormal,
			TF_OPERATOR_OBJECT_DEPLOYED,
			fmt.Sprintf("Discovered %d controller nodes", len(controllerIPs)))
	} else {
		controllerIPs = make(map[string]bool)
		for _, ip := range cr.Status.Controllers {
			// copy whatever ips available
			controllerIPs[ip] = true
		}
	}

	datapathType := DATAPATH_VPP
	if cr.Spec.UseVrouter {
		datapathType = DATAPATH_VROUTER
	}

	noLabels := []string{}
	for _, name := range nodes.WorkerNodes {
		// enable agent for all nodes
		e = SetNodeLabels(r.client, name, noLabels, datapathType)
		if e != nil {
			return false, e
		}
	}

	allLabels := []string{NODE_ROLE_ANALYTICS,
				NODE_ROLE_ANALYTICS_ALARM,
				NODE_ROLE_ANALYTICS_SNMP,
				NODE_ROLE_CONFIG,
				NODE_ROLE_CONTROL,
				NODE_ROLE_WEBUI}
	for ip, name := range nodes.MasterNodes {
		// enable all labels for master nodes
		if _, found := controllerIPs[ip]; found {
			e = SetNodeLabels(r.client, name, allLabels, datapathType)
		} else {
			e = SetNodeLabels(r.client, name, noLabels, datapathType)
		}
		if e != nil {
			return false, e
		}
	}

	var controllerNodes string
	for ip, _ := range controllerIPs {
		if controllerNodes == "" {
			controllerNodes = ip
		} else {
			controllerNodes = controllerNodes + "," + ip
		}
	}

	data := render.MakeRenderData()
	data.Data["K8S_PROVIDER"] = utils.GetKubernetesProvider()
	data.Data["TF_NAMESPACE"] = values.TF_NAMESPACE
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
	data.Data["CONTAINER_REGISTRY"] = "atsgen"
	data.Data["CONTAINER_TAG"] = cr.Spec.ReleaseTag

	if cr.Spec.ClusterName == "" {
		data.Data["KUBERNETES_CLUSTER_NAME"] = "k8s"
	} else {
		data.Data["KUBERNETES_CLUSTER_NAME"] = cr.Spec.ClusterName
	}

	updateIPForwarding(&data, cr)
	if cr.Spec.UseHostNewtorkService {
		data.Data["KUBERNETES_HOST_NETWORK_SERVICE"] = "true"
	} else {
		data.Data["KUBERNETES_HOST_NETWORK_SERVICE"] = "false"
	}

	if utils.IsOpenShiftCluster() {
		// we don't support building KMOD for openshift
		data.Data["TUNGSTEN_KMOD"] = "init"
		data.Data["CNI_BIN_DIR"] = values.OPENSHIFT_CNI_BIN_DIR
		networkConfig := &ocv1.Network{}
		err := r.client.Get(context.TODO(),
			types.NamespacedName{Name: values.OPENSHIFT_NETWORK_CONFIG,},
			networkConfig)
		if err != nil {
			log.Info("Failed to fetch openshift network config " + err.Error());
			return false, err
		}
		if (networkConfig.Spec.DisableMultiNetwork == nil ||
			!(*networkConfig.Spec.DisableMultiNetwork)) {
			data.Data["CNI_CONF_DIR"] = values.OPENSHIFT_MULTUS_CONF_DIR
		} else {
			data.Data["CNI_CONF_DIR"] = values.OPENSHIFT_CNI_CONF_DIR
		}
	} else {
		data.Data["TUNGSTEN_KMOD"] = "build"
		data.Data["CNI_BIN_DIR"] = values.DEFAULT_CNI_BIN_DIR
		data.Data["CNI_CONF_DIR"] = values.DEFAULT_CNI_CONF_DIR
	}
	data.Data["CONTROLLER_NODES"] = controllerNodes
	data.Data["CONTROL_NODES"] = controllerNodes
	data.Data["JVM_EXTRA_OPTS"] = "-Xms1g -Xmx2g"
	data.Data["KAFKA_NODES"] = controllerNodes
	data.Data["KUBERNETES_API_SECURE_PORT"] = utils.GetKubernetesApiPort()
	apiServer := utils.GetKubernetesApiServer()
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
	data.Data["VROUTER_GATEWAY"] = ""
	data.Data["WEBUI_NODES"] = controllerNodes
	data.Data["WEBUI_VIP"] = ""
	data.Data["ZOOKEEPER_PORT"] = "2181"
	data.Data["ZOOKEEPER_PORTS"] = "2888:3888"
	data.Data["DPDK_UIO_DRIVER"] = "igb_uio"
	data.Data["KUBERNETES_POD_SUBNETS"] = cr.Spec.PodNetwork.Cidr
	data.Data["KUBERNETES_SERVICE_SUBNETS"] = cr.Spec.ServiceNetwork.Cidr
	data.Data["KUBERNETES_IP_FABRIC_SUBNETS"] = cr.Spec.IpFabricNetwork.Cidr

	manifests, err := render.RenderDir(filepath.Join("/bindata", "tungsten/"), &data)
	if err != nil {
		log.Info("Failed to render yaml files " + err.Error());
		return false, err
	}

	objs = append(objs, manifests...)
	if utils.IsOpenShiftCluster() {
                // cluster is running for openshift, load objects needed for openshift
		manifests, err := render.RenderDir(filepath.Join("/bindata", "openshift/"), &data)
		if err != nil {
			log.Info("Failed to render yaml files " + err.Error());
			return false, err
		}
		objs = append(objs, manifests...)
	}

	for _, obj := range objs {
		if err := controllerutil.SetControllerReference(cr, obj, r.scheme); err!= nil {
			log.Info(err.Error())
			return false, err
		}
		if err := apply.ApplyObject(context.TODO(), r.client, obj); err != nil {
			log.Info(err.Error())
			err = errors.Wrapf(err, "could not apply (%s) %s/%s", obj.GroupVersionKind(), obj.GetNamespace(), obj.GetName())
			return false, err
		}
	}
	return true, nil
}

// Add creates a new TungstenCNI Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTungstenCNI{client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
			recorder: mgr.GetEventRecorderFor("tf-operator")}
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

	// map node scaling events to configure CNI objects as needed
	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      values.TF_OPERATOR_CONFIG,
				}},
			}
		})


	// we are intrested only in addition/removal of nodes and we will be
	// ignoring update events for nodes
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			// any addition of new node should trigger reconcile
			log.Info("Trigger reconcile, Detected new node: " + e.Meta.GetName())
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// any removal of node should trigger reconcile
			log.Info("Trigger reconcile, node removed: " + e.Meta.GetName())
			return true
		},
	}

	err = c.Watch(
		&source.Kind{Type: &corev1.Node{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: mapFn,
		},
		p)
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
	recorder record.EventRecorder
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
		r.recorder.Event(instance, corev1.EventTypeWarning,
			TF_OPERATOR_OBJECT_IGNORED, d)
		return reconcile.Result{}, nil
	}

	deployed, err := r.renderTungstenFabricCNI(instance)
	if err != nil {
		log.Error(err, "failed to reconcile")
		return reconcile.Result{}, err
	}

	if !deployed {
		s = TF_OPERATOR_OBJECT_PENDING
	}
	r.updateStatus(instance, s, d)
	log.Info("reconcile completed: Tungsten CNI " + instance.Name + " Updated")
	return reconcile.Result{}, nil
}

func (r *ReconcileTungstenCNI) updateControllerIPs(cr *tungstenv1alpha1.TungstenCNI) error {
	for ip, _ := range controllerIPs {
		cr.Status.Controllers = append(cr.Status.Controllers, ip)
	}
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		log.Error(err, "failed to update TungstenCNI status")
		return err
	}
	return nil
}

func (r *ReconcileTungstenCNI) updateStatus(cr *tungstenv1alpha1.TungstenCNI, state string, msg string) error {
        if (cr.Status.State == state && cr.Status.Error == msg &&
		cr.Status.ReleaseTag == cr.Spec.ReleaseTag) {
		// No update required
		return nil
	}
	cr.Status.ReleaseTag = cr.Spec.ReleaseTag
	cr.Status.State = state
	cr.Status.Error = msg
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		log.Error(err, "failed to update TungstenCNI status")
		return err
	}
	return nil
}
