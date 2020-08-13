package hcmanager



import (
	"testing"
	"fmt"
	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	hwcc "hardware-classification-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	"k8s.io/klog/klogr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"hardware-classification-controller/controllers"
)
/*
func benchmarkConvertBytesToGb (input bmh.Capacity,b *testing.B){
	var inGb bmh.Capacity
	for i:=0; i<b.N; i++ {
		inGb =	ConvertBytesToGb(input)
	}
	fmt.Printf("Size %d GB\n",inGb)
}
*/
func benchmarkValidateExtractedHardwareProfile (b *testing.B){
//	hardwareClassification := &hwcc.HardwareClassification{}
//	hcReconciler *HardwareClassificationReconciler
//	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration
//	hcManager := NewHardwareClassificationManager(hcReconciler.Client, hcReconciler.Log)
	hostTest := getHosts()
	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
        hcManager := NewHardwareClassificationManager(c, klogr.New())
	profile := hwcc.ExpectedHardwareConfiguration {}
/*					CPU: CPU{
						MinimumCount: 2,
						MaximumCount: 3,
						MinimumSpeed: "5",
						MaximumSpeed: "10",
					},
					Disk: Disk{
						MinimumCount:1,
						MinimumIndividualSizeGB:20,
						MaximumCount:1,
						MaximumIndividualSizeGB:20,
					},
					NIC: NIC{
						MinimumSizeGB:1,
						MaximumSizeGB:2,
					},
					RAM: RAM{
						MinimumSizeGB:6,
						MaximumSizeGB:12,
					},
				}*/
	for i:=0; i<b.N; i++{
		err := hcManager.ValidateExtractedHardwareProfile(profile)
		if err != nil {
			fmt.Println("Error",err)
		}
	}
}

func BenchMarkbenchmarkValidateExtractedHardwareProfile(b *testing.B) {benchmarkValidateExtractedHardwareProfile(b)}
//func BenchmarkConvertBytesToGb(b *testing.B)  { benchmarkConvertBytesToGb(21474825484, b) }
//func BenchmarkConvertBytesToGbNegative(b *testing.B) { benchmarkConvertBytesToGb(-21474825484, b) }
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

	return []runtime.Object{&host0}
}
