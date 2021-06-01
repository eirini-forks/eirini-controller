package controllers_test

import (
	"context"

	eiriniv1 "code.cloudfoundry.org/eirini-controller/api/v1"
	uuid "github.com/hashicorp/go-uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("LrpController", func() {
	var (
		ctx          context.Context
		lrpNamespace string
		lrpName      string
		lrp          *eiriniv1.LRP
	)

	BeforeEach(func() {
		ctx = context.Background()

		lrpName = GenerateGUID()
		lrpNamespace = "test-ns-" + GenerateGUID()

		err := k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: lrpNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		lrp = &eiriniv1.LRP{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      lrpName,
			},
			Spec: eiriniv1.LRPSpec{
				Image:  "eirini/dorini",
				DiskMB: 2,
			},
		}
	})

	It("does something", func() {
		Expect(k8sClient.Create(ctx, lrp)).To(Succeed())
		Eventually(getStatefulSetItems(ctx, lrpNamespace)).Should(HaveLen(1))
	})
})

func getStatefulSetItems(ctx context.Context, lrpNamespace string) func() ([]appsv1.StatefulSet, error) {
	return func() ([]appsv1.StatefulSet, error) {
		statefulsets := appsv1.StatefulSetList{}
		err := k8sClient.List(ctx, &statefulsets, client.InNamespace(lrpNamespace))

		return statefulsets.Items, err
	}
}

func GenerateGUID() string {
	guid, err := uuid.GenerateUUID()
	Expect(err).NotTo(HaveOccurred())

	return guid[:30]
}
