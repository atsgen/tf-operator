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
	err = c.Watch(&source.Kind{Type: &tungstenv1alpha1.SDN{}}, &handler.EnqueueRequestForOwner{
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

	if instance.Name != values.OpenShiftNetworkConfig {
		log.Info("skipping OpenShift Network Config name: " + instance.Name)
		// Return and don't requeue
		return reconcile.Result{}, nil
	}

	useTungsten := false
	if instance.Spec.NetworkType == values.OpenShiftAtsgenCni {
		useTungsten = true
	}

	// Check if Tungsten CNI installation already exists
	found := &tungstenv1alpha1.SDN{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: values.TFDefaultDeployment,}, found)
	if err != nil && errors.IsNotFound(err) {
		if !useTungsten {
			// we are not set as the CNI for OpenShift ignore
			log.Info("OpenShift is not configured to use Tungsten CNI")
			return reconcile.Result{}, nil
		}
		reqLogger.Info("Creating a new Tungsten CNI", "Name", values.TFDefaultDeployment)
		// Define a new Tungsten CNI object
		cni := newSDN(instance)

		err = r.client.Create(context.TODO(), cni)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.setNetworkStatus(instance)
		if err != nil {
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

	err = r.setNetworkStatus(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// CNI already exists - don't requeue
	reqLogger.Info("Skip reconcile: CNI already exists")
	return reconcile.Result{}, nil
}

func (r *ReconcileNetwork) setNetworkStatus(cr *configv1.Network) error {
	if cr.Status.NetworkType == values.OpenShiftAtsgenCni {
		// we don't need to update anything here
		return nil
	}

	// Get resource before updating to use in the Patch call.
	patchFrom := client.MergeFrom(cr.DeepCopy())
	cr.Status.ClusterNetwork = cr.Spec.ClusterNetwork
	cr.Status.ServiceNetwork = cr.Spec.ServiceNetwork
	cr.Status.NetworkType = values.OpenShiftAtsgenCni
	// TODO(prabhjot) for OpenShift we need to report MTU as per system
	// capabilities. However, when do VxLAN forwarding to account for tunnel
	// headers in a default environment we will be reporting 1410 as the MTU
	// to OpenShift infra. This is small value, but should work in general for
	// most of the deployments
	cr.Status.ClusterNetworkMTU = 1410
	if err := r.client.Patch(context.Background(), cr, patchFrom); err != nil {
		log.Info("Error patching openshift network status " + err.Error())
		return err
	}
	return nil
}

// newSDN returns a new tungsten CNI object
func newSDN(cr *configv1.Network) *tungstenv1alpha1.SDN {
	cni := &tungstenv1alpha1.SDN{
		ObjectMeta: metav1.ObjectMeta{
			Name:      values.TFDefaultDeployment,
		},
		Spec: tungstenv1alpha1.SDNSpec{
			ReleaseTag:    values.TFReleaseTag,
			UseVrouter:    true,
			IpForwarding:  "snat",
			UseHostNewtorkService: true,
		},
	}

	if len(cr.Spec.ClusterNetwork) != 0 {
		cni.Spec.PodNetwork = tungstenv1alpha1.PodNetworkType{
					Cidr: cr.Spec.ClusterNetwork[0].CIDR,
				}
		// TODO(prabhjot) need to see how should this actually
		// work of OpenShift
		cni.Spec.IpFabricNetwork = tungstenv1alpha1.IpFabricNetworkType{
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
