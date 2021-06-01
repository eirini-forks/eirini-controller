/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eiriniv1 "code.cloudfoundry.org/eirini-controller/api/v1"
	"code.cloudfoundry.org/eirini/k8s/reconciler"
)

// LRPReconciler reconciles a LRP object
type LRPReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	lrpsCrClient   reconciler.LRPsCrClient
	workloadClient reconciler.LRPWorkloadCLient
}

//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LRP object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *LRPReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LRPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eiriniv1.LRP{}).
		Complete(r)
}
