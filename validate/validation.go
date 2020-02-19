package validate

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

//Comparison compare the host against the profile and filter the valid host
func Comparison(hosts []bmh.BareMetalHost, profiles []hwcc.ExpectedHardwareConfiguration) {

	validHost := make(map[interface{}][]hwcc.ExpectedHardwareConfiguration)
	fmt.Println("Inside Comparison file")
	for _, host := range hosts {
		// fmt.Printf("%+v", host.Status.HardwareDetails)
		fmt.Printf("CPU Count :- %+v \n", host.Status.HardwareDetails.CPU.Count)
		fmt.Printf("Stoarage Count :- %+v \n", host.Status.HardwareDetails.Storage)
		fmt.Printf("Nics Count :- %+v \n", host.Status.HardwareDetails.NIC)
		fmt.Printf("Ram Count :- %+v \n", host.Status.HardwareDetails.RAMMebibytes)

		for _, profile := range profiles {
			if profile.MinimumCPU.Count < host.Status.HardwareDetails.CPU.Count &&
				profile.MinimumDisk.SizeBytesGB < int64(host.Status.HardwareDetails.Storage[0].SizeBytes*1024*1024) &&
				profile.MinimumNICS.NumberOfNICS < len(host.Status.HardwareDetails.NIC) &&
				profile.MinimumRAM < host.Status.HardwareDetails.RAMMebibytes {
				newHost, ok := validHost[host.Status.HardwareDetails]
				if ok {
					validHost[host.Status.HardwareDetails] = append(newHost, profile)
				} else {
					var validProfile []hwcc.ExpectedHardwareConfiguration
					validHost[host.Status.HardwareDetails] = append(validProfile, profile)
				}

			}
		}

	}

	fmt.Println("Valid host list*************************")
	fmt.Println(validHost)

}
