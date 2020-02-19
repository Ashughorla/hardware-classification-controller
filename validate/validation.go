package validate

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

//Comparison compare the host against the profile and filter the valid host
func Comparison(hosts []bmh.BareMetalHost, profiles []hwcc.ExpectedHardwareConfiguration) {

	fmt.Println("Inside Comparison file")
	for _, host := range hosts {
		// fmt.Printf("%+v", host.Status.HardwareDetails)
		fmt.Printf("CPU Count :- %+v \n", host.Status.HardwareDetails.CPU.Count)
		fmt.Printf("Stoarage Count :- %+v \n", host.Status.HardwareDetails.Storage)
		fmt.Printf("Stoarage Count :- %+v \n", host.Status.HardwareDetails.NIC)
		fmt.Printf("Stoarage Count :- %+v \n", host.Status.HardwareDetails.RAMMebibytes)

		// for _, profile := range profiles {
		// 	if profile.MinimumCPU.Count < host.Status.HardwareDetails.CPU.Count &&
		// 		profile.MinimumDisk.SizeBytesGB < host.Status.HardwareDetails.Storage.SizeBytes*1024*1024 &&
		// 		profile.MinimumNICS.NumberOfNICS < host.Status.HardwareDetails.NIC.NoOfNics &&
		// 		profile.MinimumRAM < host.Status.HardwareDetails.RAMMebibytes {
		// 		newHost, ok := validHost[host.Spec.HardwareDetails]
		// 		if ok {
		// 			validHost[host.Spec.HardwareDetails] = append(newHost, profile)
		// 		} else {
		// 			var validProfile []ironic.ExpectedHardwareConfiguration
		// 			validHost[host.Spec.HardwareDetails] = append(validProfile, profile)
		// 		}

		// 	}
		// }

	}

}
