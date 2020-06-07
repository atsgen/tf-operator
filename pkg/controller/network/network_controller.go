package network

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	tungstenv1alpha1 "github.com/atsgen/tf-operator/pkg/apis/tungsten/v1alpha1"
	"github.com/atsgen/tf-operator/pkg/values"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_network")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Network Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNetwork{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("network-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Network
	err = c.Watch(&source.Kind{Type: &configv1.Network{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Network
	err = c.Watch(&source.Kind{Type: &tungstenv1alpha1.TungstenCNI{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &configv1.Network{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNetwork implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNetwork{}

// ReconcileNetwork reconciles a Network object
type ReconcileNetwork struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Network object and makes changes based on the state read
// and what is in the Network.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNetwork) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Network")

	// Fetch the Network instance
	instance := &configv1.Network{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Name != values.OPENSHIFT_NETWORK_CONFIG {
		log.Info("skipping OpenShift Network Config name: " + instance.Name)
		// Return and don't requeue
		return reconcile.Result{}, nil
	}

	useTungsten := false
	if instance.Spec.NetworkType == values.OPENSHIFT_ATSGEN_CNI {
		useTungsten = true
	}

	// Check if Tungsten CNI installation already exists
	found := &tungstenv1alpha1.TungstenCNI{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: values.TF_OPERATOR_CONFIG,}, found)
	if err != nil && errors.IsNotFound(err) {
		if !useTungsten {
			// we are not set as the CNI for OpenShift ignore
			log.Info("OpenShift is not configured to use Tungsten CNI")
			return reconcile.Result{}, nil
		}
		reqLogger.Info("Creating a new Tungsten CNI", "Name", values.TF_OPERATOR_CONFIG)
		// Define a new Tungsten CNI object
		cni := newTungstenCNI(instance)

		err = r.client.Create(context.TODO(), cni)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Get resource before updating to use in the Patch call.
		patchFrom := client.MergeFrom(instance.DeepCopy())
		instance.Status.ClusterNetwork = instance.Spec.ClusterNetwork
		instance.Status.ServiceNetwork = instance.Spec.ServiceNetwork
		instance.Status.NetworkType = values.OPENSHIFT_ATSGEN_CNI
		// TODO(prabhjot) for OpenShift we need to report MTU as per system
		// capabilities. However, when do VxLAN forwarding to account for tunnel
		// headers in a default environment we will be reporting 1410 as the MTU
		// to OpenShift infra. This is small value, but should work in general for
		// most of the deployments
		instance.Status.ClusterNetworkMTU = 1410
		if err = r.client.Patch(context.Background(), instance, patchFrom); err != nil {
			log.Info("Error patching openshift network status", err, reqLogger.WithValues("openshiftConfig", instance))
			return reconcile.Result{}, err
		}
		// CNI created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	if !useTungsten {
		// we are not set as the CNI for OpenShift ignore
		log.Info("OpenShift is not configured to use Tungsten CNI, Deleting")
		err = r.client.Delete(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// CNI already exists - don't requeue
	reqLogger.Info("Skip reconcile: CNI already exists")
	return reconcile.Result{}, nil
}

// newTungstenCNI returns a new tungsten CNI object
func newTungstenCNI(cr *configv1.Network) *tungstenv1alpha1.TungstenCNI {
	cni := &tungstenv1alpha1.TungstenCNI{
		ObjectMeta: metav1.ObjectMeta{
			Name:      values.TF_OPERATOR_CONFIG,
		},
		Spec: tungstenv1alpha1.TungstenCNISpec{
			ReleaseTag:    values.TF_RELEASE_TAG,
			UseVrouter:    true,
		},
	}

	if len(cr.Spec.ClusterNetwork) != 0 {
		cni.Spec.PodNetwork = tungstenv1alpha1.PodNetworkType{
					Cidr: cr.Spec.ClusterNetwork[0].CIDR,
				}
	}
	if len(cr.Spec.ServiceNetwork) != 0 {
		cni.Spec.ServiceNetwork = tungstenv1alpha1.ServiceNetworkType{
					Cidr: cr.Spec.ServiceNetwork[0],
				}
	}
	return cni
}
