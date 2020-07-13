package unit

import (
	"testing"
	"fmt"
	"hardware-classification-controller/controllers"
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
	profile := hcmanager.ExpectedHardwareConfiguration {}
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
