package stset

import (
	"context"

	"code.cloudfoundry.org/eirini/api"
	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/util/retry"
)

//counterfeiter:generate . StatefulSetUpdater

type StatefulSetUpdater interface {
	Update(ctx context.Context, namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
}

type Updater struct {
	logger             lager.Logger
	statefulSetUpdater StatefulSetUpdater
	getStatefulSet     getStatefulSetFunc
	pdbUpdater         PodDisruptionBudgetUpdater
}

func NewUpdater(
	logger lager.Logger,
	statefulSetGetter StatefulSetByLRPIdentifierGetter,
	statefulSetUpdater StatefulSetUpdater,
	pdbUpdater PodDisruptionBudgetUpdater,
) Updater {
	return Updater{
		logger:             logger,
		statefulSetUpdater: statefulSetUpdater,
		pdbUpdater:         pdbUpdater,
		getStatefulSet:     newGetStatefulSetFunc(statefulSetGetter),
	}
}

func (u *Updater) Update(ctx context.Context, lrp *api.LRP) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return u.update(ctx, lrp)
	})

	return errors.Wrap(err, "failed to update statefulset")
}

func (u *Updater) update(ctx context.Context, lrp *api.LRP) error {
	logger := u.logger.Session("update", lager.Data{"guid": lrp.GUID, "version": lrp.Version})

	statefulSet, err := u.getStatefulSet(ctx, api.LRPIdentifier{GUID: lrp.GUID, Version: lrp.Version})
	if err != nil {
		logger.Error("failed-to-get-statefulset", err)

		return err
	}

	updatedStatefulSet, err := u.getUpdatedStatefulSetObj(statefulSet,
		lrp.TargetInstances,
		lrp.LastUpdated,
		lrp.Image,
	)
	if err != nil {
		logger.Error("failed-to-get-updated-statefulset", err)

		return err
	}

	if _, err = u.statefulSetUpdater.Update(ctx, updatedStatefulSet.Namespace, updatedStatefulSet); err != nil {
		logger.Error("failed-to-update-statefulset", err, lager.Data{"namespace": statefulSet.Namespace})

		return errors.Wrap(err, "failed to update statefulset")
	}

	if err = u.pdbUpdater.Update(ctx, statefulSet, lrp); err != nil {
		logger.Error("failed-to-update-disruption-budget", err, lager.Data{"namespace": statefulSet.Namespace})

		return errors.Wrap(err, "failed to delete pod disruption budget")
	}

	return nil
}

func (u *Updater) getUpdatedStatefulSetObj(sts *appsv1.StatefulSet, instances int, lastUpdated, image string) (*appsv1.StatefulSet, error) {
	updatedSts := sts.DeepCopy()

	count := int32(instances)
	updatedSts.Spec.Replicas = &count
	updatedSts.Annotations[AnnotationLastUpdated] = lastUpdated

	if image != "" {
		for i, container := range updatedSts.Spec.Template.Spec.Containers {
			if container.Name == ApplicationContainerName {
				updatedSts.Spec.Template.Spec.Containers[i].Image = image
			}
		}
	}

	return updatedSts, nil
}
