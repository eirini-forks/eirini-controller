module code.cloudfoundry.org/eirini-controller

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.20.3
	k8s.io/client-go => k8s.io/client-go v0.20.3
)

require (
	code.cloudfoundry.org/eirini v0.0.0-20210527142840-39e7adeb20ee
	github.com/hashicorp/go-uuid v1.0.2
	github.com/onsi/ginkgo v1.16.3
	github.com/onsi/gomega v1.13.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v1.5.2
	sigs.k8s.io/controller-runtime v0.8.3
)
