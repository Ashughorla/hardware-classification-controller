package controllers

import (
	"context"
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Hardware Classification Controller", func() {

	hostTest := getHosts()

	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
	r := HardwareClassificationReconciler{
		Client: c,
		Log:    klogr.New(),
	}

	It("Should Check the matched fetch host", func() {
		result := fetchBmhHostList(context.TODO(), &r, getNamespace())
		fmt.Println(result)
		Expect(len(hostTest)).Should(Equal(len(result)))
	})

})

//-----------------
// Helper functions
//-----------------
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
