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
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"code.cloudfoundry.org/eirini"
	eiriniv1 "code.cloudfoundry.org/eirini-controller/api/v1"
	"code.cloudfoundry.org/eirini/api"
	"code.cloudfoundry.org/eirini/k8s/reconciler"
	"code.cloudfoundry.org/eirini/util"
	"code.cloudfoundry.org/lager"
)

// LRPReconciler reconciles a LRP object
type LRPReconciler struct {
	client.Client
	Logger         lager.Logger
	Scheme         *runtime.Scheme
	WorkloadClient reconciler.LRPWorkloadCLient
}

//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=eirini.cloudfoundry.org,resources=lrps/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;watch;list
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=create;update

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
	// TODO maybe create a lager implementation that forwards to logr.Logger
	_ = log.FromContext(ctx)
	logger := r.Logger.Session(
		"reconcile-lrp",
		lager.Data{
			"name":      req.NamespacedName.Name,
			"namespace": req.NamespacedName.Namespace,
		},
	)

	lrp := eiriniv1.LRP{}
	if err := r.Get(ctx, req.NamespacedName, &lrp); err != nil {
		// logger.Error("failed-to-get-lrp", err)
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.do(ctx, &lrp); err != nil {
		logger.Error("failed-to-reconcile", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *LRPReconciler) do(ctx context.Context, lrp *eiriniv1.LRP) error {
	_, err := r.WorkloadClient.Get(ctx, api.LRPIdentifier{
		GUID:    lrp.Spec.GUID,
		Version: lrp.Spec.Version,
	})
	if errors.Is(err, eirini.ErrNotFound) {
		appLRP, parseErr := toAPILrp(lrp)
		if parseErr != nil {
			return errors.Wrap(parseErr, "failed to parse the crd spec to the lrp model")
		}

		return errors.Wrap(r.WorkloadClient.Desire(ctx, lrp.Namespace, appLRP, r.setOwnerFn(lrp)), "failed to desire lrp")
	}

	if err != nil {
		return errors.Wrap(err, "failed to get lrp")
	}

	appLRP, err := toAPILrp(lrp)
	if err != nil {
		return errors.Wrap(err, "failed to parse the crd spec to the lrp model")
	}

	var errs *multierror.Error

	err = r.updateStatus(ctx, lrp)
	errs = multierror.Append(errs, errors.Wrap(err, "failed to update lrp status"))

	err = r.WorkloadClient.Update(ctx, appLRP)
	errs = multierror.Append(errs, errors.Wrap(err, "failed to update app"))

	return errs.ErrorOrNil()
}

func (r *LRPReconciler) updateStatus(ctx context.Context, lrp *eiriniv1.LRP) error {
	lrpStatus, err := r.WorkloadClient.GetStatus(ctx, api.LRPIdentifier{
		GUID:    lrp.Spec.GUID,
		Version: lrp.Spec.Version,
	})
	if err != nil {
		return err
	}

	actualStaus := eiriniv1.LRPStatus{
		Replicas: lrpStatus.Replicas,
	}

	return r.UpdateLRPStatus(ctx, lrp, actualStaus)
}

func (r *LRPReconciler) UpdateLRPStatus(ctx context.Context, lrp *eiriniv1.LRP, newStatus eiriniv1.LRPStatus) error {
	newLRP := lrp.DeepCopy()
	newLRP.Status = newStatus

	return r.Status().Patch(ctx, newLRP, client.MergeFrom(lrp))
}

func (r *LRPReconciler) setOwnerFn(lrp *eiriniv1.LRP) func(interface{}) error {
	return func(resource interface{}) error {
		obj, ok := resource.(metav1.Object)
		if !ok {
			return fmt.Errorf("failed to cast %v to metav1.Object", resource)
		}

		if err := ctrl.SetControllerReference(lrp, obj, r.Scheme); err != nil {
			return errors.Wrap(err, "failed to set controller reference")
		}

		return nil
	}
}

func toAPILrp(lrp *eiriniv1.LRP) (*api.LRP, error) {
	apiLrp := &api.LRP{}
	if err := copier.Copy(apiLrp, lrp.Spec); err != nil {
		return nil, errors.Wrap(err, "failed to copy lrp spec")
	}

	apiLrp.TargetInstances = lrp.Spec.Instances

	if lrp.Spec.PrivateRegistry != nil {
		apiLrp.PrivateRegistry.Server = util.ParseImageRegistryHost(lrp.Spec.Image)
	}

	return apiLrp, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LRPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eiriniv1.LRP{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
