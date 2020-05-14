package controllers

import (
	"k8s.io/apimachinery/pkg/runtime"

	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getNamespace() string {
	return "metal3"
}

func getExtractedHardwareProfile() []runtime.Object {
	expectedHardwareClassification := hwcc.ExpectedHardwareConfiguration{
		CPU: &hwcc.CPU{
			MaximumCount: 1,
			MinimumCount: 1,
			MaximumSpeed: "1.2",
			MinimumSpeed: "1",
		},
		Disk: &hwcc.Disk{
			MaximumCount:            2,
			MaximumIndividualSizeGB: 1000,
			MinimumCount:            1,
			MinimumIndividualSizeGB: 500,
		},
		NIC: &hwcc.NIC{
			MaximumCount: 2,
			MinimumCount: 1,
		},
		RAM: &hwcc.RAM{
			MaximumSizeGB: 16,
			MinimumSizeGB: 8,
		},
	}

	expectedHardwareConfiguration := hwcc.HardwareClassification{
		Spec: hwcc.HardwareClassificationSpec{
			ExpectedHardwareConfiguration: expectedHardwareClassification,
		},
	}

	return []runtime.Object{&expectedHardwareConfiguration}
}

func getHosts() []runtime.Object {

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

	return []runtime.Object{&host0, &host1, &host2, &host3}
}
