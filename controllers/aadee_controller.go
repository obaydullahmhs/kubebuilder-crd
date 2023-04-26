/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aadeeappsv1 "github.com/obaydullahmhs/kubebuilder-crd/api/v1"
)

// AadeeReconciler reconciles a Aadee object
type AadeeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aadee.apps.obaydullahmhs,resources=aadees,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aadee.apps.obaydullahmhs,resources=aadees/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aadee.apps.obaydullahmhs,resources=aadees/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Aadee object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AadeeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.WithValues("ReqName", req.Name, "ReqNamespace", req.Namespace)
	// TODO(user): your logic here
	/*
		### 1: Load the Aadee by name

		We'll fetch the Aadee using our client.  All client methods take a
		context (to allow for cancellation) as their first argument, and the object
		in question as their last.  Get is a bit special, in that it takes a
		[`NamespacedName`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client?tab=doc#ObjectKey)
		as the middle argument (most don't have a middle argument, as we'll see
		below).

		Many client methods also take variadic options at the end.
	*/
	var aadee aadeeappsv1.Aadee

	if err := r.Get(ctx, req.NamespacedName, &aadee); err != nil {
		log.Error(err, "unable to fetch aadee")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	/*
		### 2: List all active deployments, and update the status

		To fully update our status, we'll need to list all child deployments in this namespace that belong to this Aadee.
		Similarly to Get, we can use the List method to list the child deployments.  Notice that we use variadic options to
		set the namespace and field match (which is actually an index lookup that we set up below).
	*/
	var childDeployments appsv1.DeploymentList
	if err := r.List(ctx, &childDeployments, client.InNamespace(req.Namespace), client.MatchingFields{deployOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child deployments")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AadeeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aadeeappsv1.Aadee{}).
		Complete(r)
}
