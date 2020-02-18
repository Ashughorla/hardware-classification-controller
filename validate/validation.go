package validate

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

func Validate(hosts []bmh.BareMetalHost, profiles []hwcc.ExpectedHardwareConfiguration) {

	fmt.Println("********************************")
	fmt.Println("Validation wali file")
	fmt.Printf("%+v", hosts)
	// for _, host := range hosts {
	//     fmt.Printf("%+v", host)
	//     fmt.Println("")
	//     fmt.Println("")

	// }

	// fmt.Println("*********************")

	// for _, profile := range profiles {
	//     fmt.Printf("%+v", profile)
	//     fmt.Println("")
	//     fmt.Println("")
	// }

}
