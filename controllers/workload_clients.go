package controllers

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/k8s/client"
	"code.cloudfoundry.org/eirini/k8s/pdb"
	"code.cloudfoundry.org/eirini/k8s/reconciler"
	"code.cloudfoundry.org/eirini/k8s/stset"
	"code.cloudfoundry.org/eirini/prometheus"
	"code.cloudfoundry.org/lager"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/clock"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

func CreateLRPWorkloadsClient(
	logger lager.Logger,
	controllerClient ctrlruntimeclient.Client,
	clientset kubernetes.Interface,
	cfg eirini.ControllerConfig,
	scheme *runtime.Scheme,
	latestMigration int,
) (reconciler.LRPWorkloadCLient, error) {
	logger = logger.Session("lrp-reconciler")
	lrpToStatefulSetConverter := stset.NewLRPToStatefulSetConverter(
		cfg.ApplicationServiceAccount,
		cfg.RegistrySecretName,
		cfg.UnsafeAllowAutomountServiceAccountToken,
		cfg.AllowRunImageAsRoot,
		latestMigration,
		k8s.CreateLivenessProbe,
		k8s.CreateReadinessProbe,
	)
	workloadClient := k8s.NewLRPClient(
		logger.Session("stateful-set-desirer"),
		client.NewSecret(clientset),
		client.NewStatefulSet(clientset, cfg.WorkloadsNamespace),
		client.NewPod(clientset, cfg.WorkloadsNamespace),
		pdb.NewUpdater(client.NewPodDisruptionBudget(clientset)),
		client.NewEvent(clientset),
		lrpToStatefulSetConverter,
		stset.NewStatefulSetToLRPConverter(),
	)

	return prometheus.NewLRPClientDecorator(logger.Session("prometheus-decorator"), workloadClient, metrics.Registry, clock.RealClock{})
}
