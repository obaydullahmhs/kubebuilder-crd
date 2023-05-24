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
	"fmt"
	aadeeappsv1 "github.com/obaydullahmhs/kubebuilder-crd/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
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
	logs := log.FromContext(ctx)
	logs.WithValues("ReqName", req.Name, "ReqNamespace", req.Namespace)
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
		fmt.Println(err, "unable to fetch aadee")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}
	fmt.Println("Name:", aadee.Name)

	// deploymentObject carry the all data of deployment in specific namespace and name
	var deploymentObject appsv1.Deployment
	//Naming of deployment
	deploymentName := func() string {
		name := aadee.Name + "-" + aadee.Spec.DeploymentName + "-depl"
		if aadee.Spec.DeploymentName == "" {
			// We choose to absorb the error here as the worker would requeue the
			// resource otherwise. Instead, the next time the resource is updated
			// the resource will be queued again.
			name = aadee.Name + "-missingname-depl"
		}
		return name
	}()
	//Creating NamespacedName for deploymentObject
	objectkey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}

	if err := r.Get(ctx, objectkey, &deploymentObject); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Deployment for ", aadee.Name, ", creating one...")
			err := r.Client.Create(ctx, newDeployment(&aadee, deploymentName))
			if err != nil {
				fmt.Printf("error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", aadee.Name)
			}
		} else {
			fmt.Printf("error fetching deployment %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if aadee.Spec.Replicas != nil && *aadee.Spec.Replicas != *deploymentObject.Spec.Replicas {

			fmt.Println(*aadee.Spec.Replicas, *deploymentObject.Spec.Replicas)
			fmt.Println("Deployment replica don't match.....updating")
			//As the replica count didn't match, we need to update it
			deploymentObject.Spec.Replicas = aadee.Spec.Replicas
			if err := r.Update(ctx, &deploymentObject); err != nil {
				fmt.Printf("error updating deployment %s\n", err)
				return ctrl.Result{}, err
			}
			fmt.Println("deployment updated")
		}
	}

	// here is service reconcile
	var serviceObject corev1.Service
	// service name
	serviceName := func() string {
		name := aadee.Name + "-" + aadee.Spec.Service.ServiceName + "-svc"
		if aadee.Spec.Service.ServiceName == "" {
			name = aadee.Name + "-missingname-svc"
		}
		return name
	}()
	//object key
	objectkey = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}

	// create or update service
	if err := r.Get(ctx, objectkey, &serviceObject); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Service for ", aadee.Name, ", creating one...")
			err := r.Create(ctx, newService(&aadee, serviceName))
			if err != nil {
				fmt.Printf("error while creating service %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Service Created...\n", aadee.Name)
			}
		} else {
			fmt.Printf("error fetching service %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if aadee.Spec.Replicas != nil && *aadee.Spec.Replicas != aadee.Status.AvailableReplicas {
			fmt.Println("Is this problem?")
			fmt.Println("Service replica miss match.....updating")
			aadee.Status.AvailableReplicas = *aadee.Spec.Replicas
			if err := r.Status().Update(ctx, &aadee); err != nil {
				fmt.Printf("error while updating service %s\n", err)
				return ctrl.Result{}, err
			}
			fmt.Println("service updated")
		}
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 1 * time.Minute,
	}, nil
}

var (
	deployOwnerKey = ".metadata.controller"
	svcOwnerKey    = ".metadata.controller"
	ourApiGVStr    = aadeeappsv1.GroupVersion.String()
	ourKind        = "Aadee"
)

// SetupWithManager sets up the controller with the Manager.
func (r *AadeeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Indexing our Owns resource.This will allow for quickly answer the question:
	// If Owns Resource x is updated, which CustomR are affected?
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, deployOwnerKey, func(object client.Object) []string {
		// grab the deployment object, extract the owner.
		deployment := object.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
		if owner == nil {
			return nil
		}
		// make sure it's a Aadee
		if owner.APIVersion != ourApiGVStr || owner.Kind != ourKind {
			return nil
		}
		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Service{}, svcOwnerKey, func(object client.Object) []string {
		svc := object.(*corev1.Service)
		owner := metav1.GetControllerOf(svc)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != ourApiGVStr || owner.Kind != ourKind {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	// Implementation with watches and custom eventHandler
	// if someone edit the resources(here example given for deployment resource) by kubectl
	handlerForDeployment := handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
		// List all the CR
		aadees := &aadeeappsv1.AadeeList{}
		if err := r.List(context.Background(), aadees); err != nil {
			return nil
		}
		// This func return a reconcile request array
		var request []reconcile.Request
		for _, deploy := range aadees.Items {
			deploymentName := func() string {
				name := deploy.Name + "-" + deploy.Spec.DeploymentName + "-depl"
				if deploy.Spec.DeploymentName == "" {
					name = deploy.Name + "-missingname-depl"
				}
				return name
			}()
			// Find the deployment owned by the CR
			if deploymentName == obj.GetName() && deploy.Namespace == obj.GetNamespace() {
				curdeploy := &appsv1.Deployment{}
				if err := r.Get(context.Background(), types.NamespacedName{
					Namespace: obj.GetNamespace(),
					Name:      obj.GetName(),
				}, curdeploy); err != nil {
					// This case can happen if somehow deployment gets deleted by
					// Kubectl command. We need to append new reconcile request to array
					// to create desired number of deployment again.
					if errors.IsNotFound(err) {
						request = append(request, reconcile.Request{
							NamespacedName: types.NamespacedName{
								Namespace: deploy.Namespace,
								Name:      deploy.Name,
							},
						})
						continue
					} else {
						return nil
					}
				}
				// Only append to the reconcile request array if replica count miss match.
				if curdeploy.Spec.Replicas != deploy.Spec.Replicas {
					request = append(request, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Namespace: deploy.Namespace,
							Name:      deploy.Name,
						},
					})
				}
			}
		}
		return request
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&aadeeappsv1.Aadee{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, handlerForDeployment).
		Owns(&corev1.Service{}).
		Complete(r)
	//for owns only
	//return ctrl.NewControllerManagedBy(mgr).
	//	For(&aadeeappsv1.Aadee{}).
	//	Owns(&appsv1.Deployment{}).
	//	Owns(&corev1.Service{}).
	//	Complete(r)
}
