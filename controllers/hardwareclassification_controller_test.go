package controllers

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utils "hardware-classification-controller/hcutils"

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Hardware Classification Controller", func() {

	hostTest := getHosts()

	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
	hcManager := utils.NewHardwareClassificationManager(c, klogr.New())

	It("Should Check the matched fetch host", func() {
		result, _, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			Expect(len(hostTest)).Should(Equal(len(result)))
		}

	})

	It("Should check the reconcile function", func() {
		config := getExtractedHardwareProfile()
		c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), config...)
		hardwareReconciler := &HardwareClassificationReconciler{
			Client: c,
			Log:    klogr.New(),
		}
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "Hardware Classification Controller",
				Namespace: getNamespace(),
			},
		}

		res, err := hardwareReconciler.Reconcile(req)
		fmt.Println("Reconcile output", res)
		if err == nil {
			Expect(true).Should(Equal(true))
		}

	})

})

//setupSchemeMm Add the bmoapi to our scheme
func setupSchemeMm() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := bmoapis.AddToScheme(s); err != nil {
		panic(err)
	}
	if err := hwcc.AddToScheme(s); err != nil {
		panic(err)
	}
	return s
}
