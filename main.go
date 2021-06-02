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

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"code.cloudfoundry.org/eirini"
	eiriniv1 "code.cloudfoundry.org/eirini-controller/api/v1"
	"code.cloudfoundry.org/eirini-controller/controllers"
	"code.cloudfoundry.org/eirini/migrations"
	"code.cloudfoundry.org/lager"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(eiriniv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	kubeconfig := ctrl.GetConfigOrDie()

	mgr, err := ctrl.NewManager(kubeconfig, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "825b0a36.cloudfoundry.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	logger := lager.NewLogger("eirini-controller")
	clientset := kubernetes.NewForConfigOrDie(kubeconfig)

	lrpWorkloadsClient, err := controllers.CreateLRPWorkloadsClient(
		logger,
		mgr.GetClient(),
		clientset,
		eirini.ControllerConfig{},
		mgr.GetScheme(),
		getLatestMigrationIndex(),
	)
	if err != nil {
		setupLog.Error(err, "unable to create lrp workloads client")
		os.Exit(1)
	}

	if true {
		panic("foo")
	}

	if err = (&controllers.LRPReconciler{
		Logger:         logger,
		Client:         mgr.GetClient(),
		Scheme:         mgr.GetScheme(),
		WorkloadClient: lrpWorkloadsClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "LRP")
		os.Exit(1)
	}
	if err = (&controllers.TaskReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Task")
		os.Exit(1)
	}
	if err = (&eiriniv1.LRP{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "LRP")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func getLatestMigrationIndex() int {
	return migrations.CreateMigrationStepsProvider(nil, nil, nil, "").GetLatestMigrationIndex()
}
