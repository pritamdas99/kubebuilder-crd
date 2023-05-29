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

package controller

import (
	"context"
	"github.com/PritamDas17021999/kubebuilder-crd/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	logger "log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"time"
)

// PritamReconciler reconciles a Pritam object
type PritamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func buildSlice(arr ...string) []string {
	return arr
}

//+kubebuilder:rbac:groups=pritamdas.dev.pritamdas.dev,resources=pritams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pritamdas.dev.pritamdas.dev,resources=pritams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pritamdas.dev.pritamdas.dev,resources=pritams/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pritam object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *PritamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx)
	logs.WithValues("ReqName", req.Name, "ReqNamesapce", req.Namespace)

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

	var pritam v1alpha1.Pritam

	if err := r.Get(ctx, req.NamespacedName, &pritam); err != nil {
		logger.Println(err, "unable to fetch pritam")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}

	logger.Println("Pritam Name", pritam.Name)

	// deploymentObject carry the all data of deployment in specific namespace and name
	var deploymentObject appsv1.Deployment

	if pritam.Spec.Name == "" {
		//	pritam_deep := pritam.DeepCopy()
		pritam.Spec.Name = pritam.Name
		r.Update(ctx, &pritam)
		logger.Println("updated pritam", pritam.Spec.Name, "again", pritam.Name)
	}

	objectKey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      pritam.Spec.Name,
	}

	if err := r.Get(ctx, objectKey, &deploymentObject); err != nil {
		if errors.IsNotFound(err) {
			logger.Println("could not find existing Deployment for ", pritam.Name, " creating one...")
			err := r.Create(ctx, newDeployment(&pritam))
			if err != nil {
				logger.Printf("error while creteing depluoyment %s\n", err)
				return ctrl.Result{}, err
			}
			logger.Printf("%s deployment created...\n", pritam.Name)

		}
		logger.Printf("error fetchiong deploymewnt %s\n", err)
		return ctrl.Result{}, err
	} else {
		if pritam.Spec.Replicas != nil && *pritam.Spec.Replicas != *deploymentObject.Spec.Replicas {
			logger.Println(*pritam.Spec.Replicas, *deploymentObject.Spec.Replicas)
			logger.Println("deployemnts replicas don't match...")
			*deploymentObject.Spec.Replicas = *pritam.Spec.Replicas
			if err := r.Update(ctx, &deploymentObject); err != nil {
				logger.Printf("error updating deployment %s\n", err)
				return ctrl.Result{}, err
			}
			logger.Println("deployment updated")
		}
		if pritam.Spec.Replicas != nil && *pritam.Spec.Replicas != pritam.Status.AvilableReplicas {
			logger.Println(*pritam.Spec.Replicas, pritam.Status.AvilableReplicas)
			var pritam_deep *v1alpha1.Pritam
			pritam_deep = pritam.DeepCopy()
			logger.Printf("Is it problem?\nService replica missmatch...")
			pritam_deep.Status.AvilableReplicas = *pritam.Spec.Replicas
			if err := r.Status().Update(ctx, pritam_deep); err != nil {
				logger.Printf("error updating status %s\n", err)
				return ctrl.Result{}, err
			}
			logger.Println("status updated")
		}
	}

	var serviceObject corev1.Service
	objectKey = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      strings.Join(buildSlice(pritam.Spec.Name, v1alpha1.Service), "-"),
	}

	if err := r.Get(ctx, objectKey, &serviceObject); err != nil {
		if errors.IsNotFound(err) {
			logger.Println("could not find existing Service for ", pritam.Name, " creating one...")
			err := r.Create(ctx, newService(&pritam))
			if err != nil {
				logger.Printf("error while creating Service %s\n", err)
				return ctrl.Result{}, err
			}
			logger.Printf("%s Service created...\n", pritam.Name)

		}
		//logger.Printf("error fetchiong service %s\n", err)
		//return ctrl.Result{}, err

	}
	logger.Printf("service updated")

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 1 * time.Minute,
	}, nil
}

var (
	deployOwnerKey = ".metadata.controller"
	svcOwnerKey    = ".metadata.controller"
	ourApiGVStr    = v1alpha1.GroupVersion.String()
	ourKind        = "Pritam"
)

// SetupWithManager sets up the controller with the Manager.
func (r *PritamReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, deployOwnerKey, func(object client.Object) []string {
		deployment := object.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
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

	handlerForDeployment := handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {

		pritams := &v1alpha1.PritamList{}
		if err := r.List(context.Background(), pritams); err != nil {
			return nil
		}
		var request []reconcile.Request
		for _, deploy := range pritams.Items {
			deploymentName := func() string {
				name := deploy.Spec.Name
				if deploy.Spec.Name == "" {
					name = strings.Join(buildSlice(deploy.Name, v1alpha1.Deployment), "-")
				}
				return name
			}()

			if deploymentName == obj.GetName() && deploy.Namespace == obj.GetNamespace() {
				dummy := &appsv1.Deployment{}
				if err := r.Get(context.Background(), types.NamespacedName{
					Namespace: obj.GetNamespace(),
					Name:      obj.GetName(),
				}, dummy); err != nil {

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

				if dummy.Spec.Replicas != deploy.Spec.Replicas {
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
		For(&v1alpha1.Pritam{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, handlerForDeployment).
		//	Watches(&source.Kind{Type: &appsv1.Deployment{}}, handlerForDeployment).
		Owns(&corev1.Service{}).
		Complete(r)
}
