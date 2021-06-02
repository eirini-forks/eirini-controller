module code.cloudfoundry.org/eirini-controller

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.20.3
	k8s.io/client-go => k8s.io/client-go v0.20.3
)

require (
	code.cloudfoundry.org/eirini v0.0.0-20210527142840-39e7adeb20ee
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-uuid v1.0.2
	github.com/jinzhu/copier v0.3.0
	github.com/onsi/ginkgo v1.16.3
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v1.5.2
	k8s.io/utils v0.0.0-20210521133846-da695404a2bc
	sigs.k8s.io/controller-runtime v0.8.3
)
