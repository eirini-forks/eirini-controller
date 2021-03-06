package stset

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/api"
	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
)

//counterfeiter:generate . StatefulSetDeleter
//counterfeiter:generate . PodDeleter

type StatefulSetDeleter interface {
	Delete(ctx context.Context, namespace string, name string) error
}

type PodDeleter interface {
	Delete(ctx context.Context, namespace, name string) error
}

type Stopper struct {
	logger             lager.Logger
	statefulSetDeleter StatefulSetDeleter
	podDeleter         PodDeleter
	getStatefulSet     getStatefulSetFunc
}

func NewStopper(
	logger lager.Logger,
	statefulSetGetter StatefulSetByLRPIdentifierGetter,
	statefulSetDeleter StatefulSetDeleter,
	podDeleter PodDeleter,
) Stopper {
	return Stopper{
		logger:             logger,
		statefulSetDeleter: statefulSetDeleter,
		podDeleter:         podDeleter,
		getStatefulSet:     newGetStatefulSetFunc(statefulSetGetter),
	}
}

func (s *Stopper) Stop(ctx context.Context, identifier api.LRPIdentifier) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return s.stop(ctx, identifier)
	})

	return errors.Wrap(err, "failed to delete statefulset")
}

func (s *Stopper) stop(ctx context.Context, identifier api.LRPIdentifier) error {
	logger := s.logger.Session("stop", lager.Data{"guid": identifier.GUID, "version": identifier.Version})
	statefulSet, err := s.getStatefulSet(ctx, identifier)

	if errors.Is(err, eirini.ErrNotFound) {
		logger.Debug("statefulset-does-not-exist")

		return nil
	}

	if err != nil {
		logger.Error("failed-to-get-statefulset", err)

		return err
	}

	if err := s.statefulSetDeleter.Delete(ctx, statefulSet.Namespace, statefulSet.Name); err != nil {
		logger.Error("failed-to-delete-statefulset", err)

		return errors.Wrap(err, "failed to delete statefulset")
	}

	return nil
}

func (s *Stopper) StopInstance(ctx context.Context, identifier api.LRPIdentifier, index uint) error {
	logger := s.logger.Session("stopInstance", lager.Data{"guid": identifier.GUID, "version": identifier.Version, "index": index})
	statefulset, err := s.getStatefulSet(ctx, identifier)

	if errors.Is(err, eirini.ErrNotFound) {
		logger.Debug("statefulset-does-not-exist")

		return nil
	}

	if err != nil {
		logger.Debug("failed-to-get-statefulset", lager.Data{"error": err.Error()})

		return err
	}

	if int32(index) >= *statefulset.Spec.Replicas {
		return eirini.ErrInvalidInstanceIndex
	}

	err = s.podDeleter.Delete(ctx, statefulset.Namespace, fmt.Sprintf("%s-%d", statefulset.Name, index))
	if k8serrors.IsNotFound(err) {
		logger.Debug("pod-does-not-exist")

		return nil
	}

	return errors.Wrap(err, "failed to delete pod")
}
