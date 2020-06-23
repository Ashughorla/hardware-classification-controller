package hcmanager

import (
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

// //RAM contains ram details fetched from the introspection data
// type RAM struct {
// 	RAMGb int `json:"ramMebibytes"`
// }

// //HardwareSystemVendor contains hardware manufacturer details fetched from the introspection data
// type HardwareSystemVendor struct {
// 	Manufacturer string `json:"manufacturer"`
// }

//NICDetails contains the nic details fetched from the introspection data
type NICDetails struct {
	Nic   bmh.NIC `json:"nic"`
	Count int     `json:"count"`
}

//StorageDetails contains disk details fetched from the introspection data
type StorageDetails struct {
	Count int           `json:"count"`
	Disk  []bmh.Storage `json:"disk"`
}

// //Disk contains disk size fetched from the introspection data
// type Disk struct {
// 	Name   string `json:"name"`
// 	SizeGb int64  `json:"sizeBytes"`
// }

// //CPU contains the clockspeed and count details fetched from the introspection data
// type CPU struct {
// 	Count      int     `json:"count"`
// 	ClockSpeed float64 `json:"clockspeed"`
// }
