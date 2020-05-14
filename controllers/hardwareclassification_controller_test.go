package controllers

import (
	"context"
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Hardware Classification Controller", func() {

	fmt.Println("inside test function")
	host0 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-0",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}

	host1 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-1",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}

	host2 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-2",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}

	host3 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-3",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}

	// expectedHardwareClassification := hwcc.ExpectedHardwareConfiguration{
	// 	CPU: &hwcc.CPU{
	// 		MaximumCount: 1,
	// 		MinimumCount: 1,
	// 		MaximumSpeed: "1.2",
	// 		MinimumSpeed: "1",
	// 	},
	// }

	type testCaseForBMHostList struct {
		Hosts         []runtime.Object
		ExpectedHosts []bmh.BareMetalHost
		Namespace     string
	}

	hostTest := []runtime.Object{&host0, &host1, &host2, &host3}

	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
	r := HardwareClassificationReconciler{
		Client: c,
		Log:    klogr.New(),
	}
	result := fetchBmhHostList(context.TODO(), &r, "default")

	It("Should Check the matched fetch host", func() {
		Expect(len(result)).Should(Equal(len(hostTest)))
	})

	DescribeTable("Test Fetch Baremetal Host List",
		func(tc testCaseForBMHostList) {

			fmt.Println("Inside Describe table")

			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), tc.Hosts...)
			r := HardwareClassificationReconciler{
				Client: c,
				Log:    klogr.New(),
			}

			result := fetchBmhHostList(context.TODO(), &r, tc.Namespace)

			if len(result) != 0 {
				for i, host := range tc.ExpectedHosts {
					Expect(result[i].Name).To(Equal(host.Name))
				}

			} else {
				Fail("Unable to fetch host list")
				Expect(len(result)).To(Equal(0))
			}

		})

	Entry("Get Host from the ready state of namespace metal3", testCaseForBMHostList{
		Hosts:         []runtime.Object{&host0, &host1, &host2, &host3},
		Namespace:     "metal3",
		ExpectedHosts: []bmh.BareMetalHost{host2, host3},
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
