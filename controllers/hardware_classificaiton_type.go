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
			HardwareDetails: &bmh.HardwareDetails{
				CPU:      bmh.CPU{Arch: "x86_64", Model: "Intel(R) Xeon(R) Gold 6226 CPU @ 2.70GHz", Count: 48, ClockMegahertz: 3700},
				Firmware: bmh.Firmware{BIOS: bmh.BIOS{Date: "", Vendor: "", Version: ""}},
				Hostname: "localhost.localdomain",
				NIC: []bmh.NIC{{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth11", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.168.121.96", MAC: "b8:59:9f:cf:fa:b2", Model: "0x15b3 0x1015", Name: "eth10", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.168.121.65", MAC: "b8:59:9f:cf:fa:ba", Model: "0x15b3 0x1015", Name: "eth6", PXE: true, SpeedGbps: 0, VLANID: 0}},
				RAMMebibytes: 196608,
				Storage: []bmh.Storage{{Name: "/dev/sda", SizeBytes: 599550590976},
					{Name: "/dev/sdb", SizeBytes: 599550590976},
					{Name: "/dev/sdc", SizeBytes: 599550590976},
					{Name: "/dev/sdd", SizeBytes: 599550590976},
					{Name: "/dev/sde", SizeBytes: 599550590976},
					{Name: "/dev/sdf", SizeBytes: 599550590976},
					{Name: "/dev/sdg", SizeBytes: 599550590976},
					{Name: "/dev/sdh", SizeBytes: 599550590976},
					{Name: "/dev/sdi", SizeBytes: 599550590976}},
				SystemVendor: bmh.HardwareSystemVendor{Manufacturer: "Dell Inc.", ProductName: "PowerEdge R740xd (SKU=NotProvided;ModelName=PowerEdge R740xd)", SerialNumber: "D2XKS13"},
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
			}, HardwareDetails: &bmh.HardwareDetails{
				CPU:      bmh.CPU{Arch: "x86_64", Model: "Intel(R) Xeon(R) Gold 6226 CPU @ 2.70GHz", Count: 40, ClockMegahertz: 3400},
				Firmware: bmh.Firmware{BIOS: bmh.BIOS{Date: "", Vendor: "", Version: ""}},
				Hostname: "localhost.localdomain",
				NIC: []bmh.NIC{{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth11", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.168.121.65", MAC: "b8:59:9f:cf:fa:ba", Model: "0x15b3 0x1015", Name: "eth6", PXE: true, SpeedGbps: 0, VLANID: 0}},
				RAMMebibytes: 196608,
				Storage: []bmh.Storage{{Name: "/dev/sda", SizeBytes: 599550590976},
					{Name: "/dev/sdb", SizeBytes: 599550590976},
					{Name: "/dev/sdc", SizeBytes: 599550590976},
					{Name: "/dev/sdd", SizeBytes: 599550590976},
					{Name: "/dev/sde", SizeBytes: 599550590976},
					{Name: "/dev/sdf", SizeBytes: 599550590976},
					{Name: "/dev/sdg", SizeBytes: 599550590976},
					{Name: "/dev/sdh", SizeBytes: 599550590976},
					{Name: "/dev/sdi", SizeBytes: 599550590976}},
				SystemVendor: bmh.HardwareSystemVendor{Manufacturer: "Dell Inc.", ProductName: "PowerEdge R740xd (SKU=NotProvided;ModelName=PowerEdge R740xd)", SerialNumber: "D2XKS13"},
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
			}, HardwareDetails: &bmh.HardwareDetails{
				CPU:      bmh.CPU{Arch: "x86_64", Model: "Intel(R) Xeon(R) Gold 6226 CPU @ 2.70GHz", Count: 32, ClockMegahertz: 4400},
				Firmware: bmh.Firmware{BIOS: bmh.BIOS{Date: "", Vendor: "", Version: ""}},
				Hostname: "localhost.localdomain",
				NIC: []bmh.NIC{{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth11", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth10", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth09", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.168.121.65", MAC: "b8:59:9f:cf:fa:ba", Model: "0x15b3 0x1015", Name: "eth6", PXE: true, SpeedGbps: 0, VLANID: 0}},
				RAMMebibytes: 196608,
				Storage: []bmh.Storage{{Name: "/dev/sda", SizeBytes: 599550590976},
					{Name: "/dev/sdb", SizeBytes: 599550590976},
					{Name: "/dev/sdc", SizeBytes: 599550590976},
					{Name: "/dev/sdd", SizeBytes: 599550590976},
					{Name: "/dev/sde", SizeBytes: 599550590976},
					{Name: "/dev/sdf", SizeBytes: 599550590976},
					{Name: "/dev/sdg", SizeBytes: 599550590976},
					{Name: "/dev/sdh", SizeBytes: 599550590976},
					{Name: "/dev/sdi", SizeBytes: 599550590976}},
				SystemVendor: bmh.HardwareSystemVendor{Manufacturer: "Dell Inc.", ProductName: "PowerEdge R740xd (SKU=NotProvided;ModelName=PowerEdge R740xd)", SerialNumber: "D2XKS13"},
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
			}, HardwareDetails: &bmh.HardwareDetails{
				CPU:      bmh.CPU{Arch: "x86_64", Model: "Intel(R) Xeon(R) Gold 6226 CPU @ 2.70GHz", Count: 32, ClockMegahertz: 4400},
				Firmware: bmh.Firmware{BIOS: bmh.BIOS{Date: "", Vendor: "", Version: ""}},
				Hostname: "localhost.localdomain",
				NIC: []bmh.NIC{{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth11", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth10", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "", MAC: "b8:59:9f:cf:fa:b3", Model: "0x15b3 0x1015", Name: "eth09", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.168.121.65", MAC: "b8:59:9f:cf:fa:ba", Model: "0x15b3 0x1015", Name: "eth6", PXE: true, SpeedGbps: 0, VLANID: 0}},
				RAMMebibytes: 196608,
				Storage: []bmh.Storage{{Name: "/dev/sda", SizeBytes: 599550590976},
					{Name: "/dev/sdb", SizeBytes: 599550590976},
					{Name: "/dev/sdc", SizeBytes: 599550590976},
					{Name: "/dev/sdd", SizeBytes: 599550590976},
					{Name: "/dev/sde", SizeBytes: 599550590976},
					{Name: "/dev/sdf", SizeBytes: 599550590976},
					{Name: "/dev/sdg", SizeBytes: 599550590976},
					{Name: "/dev/sdh", SizeBytes: 599550590976},
					{Name: "/dev/sdi", SizeBytes: 599550590976}},
				SystemVendor: bmh.HardwareSystemVendor{Manufacturer: "Dell Inc.", ProductName: "PowerEdge R740xd (SKU=NotProvided;ModelName=PowerEdge R740xd)", SerialNumber: "D2XKS13"},
			},
		},
	}

	return []runtime.Object{&host0, &host1, &host2, &host3}
}
