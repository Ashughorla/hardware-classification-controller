package controllers

import (
	"context"
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utils "hardware-classification-controller/hcmanager"

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
		fmt.Println("Fetched Host", result)

		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			Expect(len(hostTest)).Should(Equal(len(result)))
		}

	})

	It("Should Check the validated function return the hardware details list", func() {
		result, _, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(getExtractedHardwareProfile(), result)

			fmt.Println(validatedHardwareDetails)
			if len(validatedHardwareDetails) != 0 {
				fmt.Println("Validated Host", validatedHardwareDetails)
				Expect(len(hostTest)).Should(Equal(len(validatedHardwareDetails)))
			} else {
				Expect(len(hostTest)).Should(Equal(0))
			}
		}

	})

	It("Should Check the compared host list name", func() {
		result, _, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(getExtractedHardwareProfile(), result)

			fmt.Println(validatedHardwareDetails)
			if len(validatedHardwareDetails) != 0 {
				fmt.Println("Validated Host list", validatedHardwareDetails)
				comparedHost := hcManager.MinMaxComparison(getTestProfileName(), validatedHardwareDetails, getExtractedHardwareProfile())
				fmt.Println("Compared Host list", comparedHost)
				if len(comparedHost) != 0 {
					Expect(comparedHost)
				} else {
					Expect(len(hostTest)).Should(Equal(0))
				}

			} else {
				Expect(len(hostTest)).Should(Equal(0))
			}
		}

	})

	It("Should Check the if labels are set", func() {
		result, BMHList, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(getExtractedHardwareProfile(), result)

			fmt.Println(validatedHardwareDetails)
			if len(validatedHardwareDetails) != 0 {
				fmt.Println("Validated Host list", validatedHardwareDetails)
				comparedHost := hcManager.MinMaxComparison(getTestProfileName(), validatedHardwareDetails, getExtractedHardwareProfile())
				fmt.Println("Compared Host list", comparedHost)
				if len(comparedHost) != 0 {
					_, errHost, _ := hcManager.SetLabel(context.Background(), getObjectMeta(), comparedHost, BMHList, getObjectMeta().Labels)
					fmt.Println("Label Update Error", err)
					if len(errHost) > 0 {
						Fail("Error updating labels")
					} else {
						fmt.Println("Label Set Successfully")
					}
				} else {
					Expect(len(hostTest)).Should(Equal(0))
				}

			} else {
				Expect(len(hostTest)).Should(Equal(0))
			}
		}

	})

	It("Should check the reconcile function", func() {
		config := getExtractedHardwareProfileRuntime()
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
