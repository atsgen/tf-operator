package sdn

import (
	"context"
	"fmt"
	"time"

	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"k8s.io/client-go/tools/record"

	ocv1 "github.com/openshift/api/operator/v1"

	"github.com/atsgen/tf-operator/pkg/utils"
	"github.com/atsgen/tf-operator/pkg/values"
)

var log = logf.Log.WithName("controller_sdn")

var controllerIPs map[string]bool = make(map[string]bool)

func (r *ReconcileSDN) updateOpenShiftMultusStatus() error {
	if !utils.IsOpenShiftCluster() {
		return nil
	}

	_, err := utils.IsOpenShiftMultusEnabled()
	if err == nil {
		// we have already executed the block below, skip re doing it
		return nil
	}

	networkConfig := &ocv1.Network{}
	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: values.OpenShiftNetworkConfig,},
		networkConfig)
	if err != nil {
		log.Info("Failed to fetch openshift network config " + err.Error());
		return err
	}
	if (networkConfig.Spec.DisableMultiNetwork == nil ||
		!(*networkConfig.Spec.DisableMultiNetwork)) {
		utils.SetOpenShiftMultusStatus(true)
	} else {
		utils.SetOpenShiftMultusStatus(false)
	}
	return nil
}

func getNumControllerIP() int {
	if utils.IsTungstenFabricHADisabled() {
		return 1
	}
	// we allow only 3 controller nodes for now
	// TODO(prabhjot) will need to consider making this
	// configurable/dynamic, possibly based on number of
	// replicas under install-config:controlPlane in
	// configmap kube-system:cluster-config-v1
	return 3
}

func (r *ReconcileSDN) generateControllerIPList(cr *tungstenv1alpha1.SDN, nodes *NodeList) (string, error) {
	if len(nodes.MasterNodes) < getNumControllerIP() {
		// return from here we will get notified when a new
		// node is available
		r.recorder.Event(cr, corev1.EventTypeNormal,
			TFOperatorObjectPending,
			fmt.Sprintf("waiting for master node discovery (got %d/%d)", len(nodes.MasterNodes), getNumControllerIP()))
		return TFOperatorObjectPending, nil
	}
	i := 0
	for ip, _ := range nodes.MasterNodes {
		controllerIPs[ip] = true
		i++
		if i == getNumControllerIP() {
			break
		}
	}
	// commit controller ips to status before going any further
	err := r.updateControllerIPs(cr)
	if err != nil {
		return "", err
	}

	r.recorder.Event(cr, corev1.EventTypeNormal,
		TFOperatorObjectDeployed,
		fmt.Sprintf("Discovered %d controller nodes", len(controllerIPs)))
	return "", nil
}

func (r *ReconcileSDN) deployTungstenFabric(cr *tungstenv1alpha1.SDN) (string, error) {
	nodes, e := FetchNodeList(r.client)

	if e != nil {
		return "", e
	}

	// check if we have already identified Controller IPs
	if len(cr.Status.Controllers) == 0 {
		str, err := r.generateControllerIPList(cr, nodes)
		if str != "" || err != nil {
			return str, err
		}
	} else {
		controllerIPs = make(map[string]bool)
		for _, ip := range cr.Status.Controllers {
			// copy whatever ips available
			controllerIPs[ip] = true
		}
	}

	datapathType := DatapathVpp
	if cr.Spec.DatapathConfig.UseVrouter {
		datapathType = DatapathVrouter
	}

	noLabels := []string{}
	for _, name := range nodes.WorkerNodes {
		// enable agent for all nodes
		e = SetNodeLabels(r.client, name, noLabels, datapathType)
		if e != nil {
			return "", e
		}
	}

	allLabels := []string{NodeRoleAnalytics,
				NodeRoleAnalyticsAlarm,
				NodeRoleAnalyticsSnmp,
				NodeRoleConfig,
				NodeRoleControl,
				NodeRoleWebui}
	for ip, name := range nodes.MasterNodes {
		// enable all labels for master nodes
		if _, found := controllerIPs[ip]; found {
			e = SetNodeLabels(r.client, name, allLabels, datapathType)
		} else {
			e = SetNodeLabels(r.client, name, noLabels, datapathType)
		}
		if e != nil {
			return "", e
		}
	}

	e = r.updateOpenShiftMultusStatus()
	if e != nil {
		return "", e
	}

	return renderTungstenFabric(r, cr, nodes)
}

// Add creates a new SDN Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSDN{client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
			recorder: mgr.GetEventRecorderFor("tf-operator")}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sdn-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SDN
	err = c.Watch(&source.Kind{Type: &tungstenv1alpha1.SDN{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// map node scaling events to configure CNI objects as needed
	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      values.TFDefaultDeployment,
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
	// Watch for changes to secondary resource Pods and requeue the owner SDN
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tungstenv1alpha1.SDN{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSDN implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSDN{}

// ReconcileSDN reconciles a SDN object
type ReconcileSDN struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a SDN object and makes changes based on the state read
// and what is in the SDN.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSDN) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the SDN instance
	instance := &tungstenv1alpha1.SDN{}
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
	if s == TFOperatorObjectIgnored {
		log.Info("Error!!! Ignoring tf-operator " + request.Name)
		r.updateStatus(instance, s, d)
		r.recorder.Event(instance, corev1.EventTypeWarning,
			TFOperatorObjectIgnored, d)
		return reconcile.Result{}, nil
	}

	s, err = r.deployTungstenFabric(instance)
	if err != nil {
		log.Error(err, "failed to reconcile")
		return reconcile.Result{}, err
	}

	r.updateStatus(instance, s, d)
	if s == TFOperatorObjectUpdating {
		// we are still in creating stage, reconcile after 15 secs
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}
	log.Info("reconcile completed: Tungsten CNI " + instance.Name + " Updated")
	return reconcile.Result{}, nil
}

func (r *ReconcileSDN) updateControllerIPs(cr *tungstenv1alpha1.SDN) error {
	for ip, _ := range controllerIPs {
		cr.Status.Controllers = append(cr.Status.Controllers, ip)
	}
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		log.Error(err, "failed to update SDN status")
		return err
	}
	return nil
}

func (r *ReconcileSDN) updateStage(crOld *tungstenv1alpha1.SDN, stage string) error {
	cr := &tungstenv1alpha1.SDN{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{Namespace: crOld.Namespace, Name: crOld.Name,},
		cr)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return err
		}
	}
	if (cr.Status.Stage == stage) {
		// No update required
		return nil
	}
	cr.Status.Stage = stage
	err = r.client.Status().Update(context.Background(), cr)
	if err != nil {
		log.Error(err, "failed to update SDN stage")
		return err
	}
	return nil
}

func (r *ReconcileSDN) updateStatus(crOld *tungstenv1alpha1.SDN, state string, msg string) error {
	cr := &tungstenv1alpha1.SDN{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{Namespace: crOld.Namespace, Name: crOld.Name,},
		cr)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return err
		}
	}
        if (cr.Status.State == state && cr.Status.Error == msg &&
		cr.Status.ReleaseTag == cr.Spec.ReleaseTag) {
		// No update required
		return nil
	}
	cr.Status.ReleaseTag = utils.GetReleaseTag(cr.Spec.ReleaseTag)
	cr.Status.State = state
	cr.Status.Error = msg
	err = r.client.Status().Update(context.Background(), cr)
	if err != nil {
		log.Error(err, "failed to update SDN status")
		return err
	}
	return nil
}
