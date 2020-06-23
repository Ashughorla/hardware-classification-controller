package hcmanager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

// ExtractAndValidateHardwareDetails this function will return map containing introspection details for a host.
func (mgr HardwareClassificationManager) ExtractAndValidateHardwareDetails(extractedProfile hwcc.ExpectedHardwareConfiguration,
	bmhList []bmh.BareMetalHost) map[string]map[string]interface{} {

	validatedHostMap := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range bmhList {
			hardwareDetails := make(map[string]interface{})

			if extractedProfile.CPU != nil {
				// Get the CPU details from the baremetal host and validate it into new structure
				validCPU := bmh.CPU{
					Count:          host.Status.HardwareDetails.CPU.Count,
					ClockMegahertz: bmh.ClockSpeed(host.Status.HardwareDetails.CPU.ClockMegahertz) / 1000,
				}
				hardwareDetails[CPULabel] = validCPU
			}

			if extractedProfile.Disk != nil {
				// Get the Storage details from the baremetal host and validate it into new structure
				var disks []bmh.Storage

				for _, disk := range host.Status.HardwareDetails.Storage {
					disks = append(disks, bmh.Storage{Name: disk.Name, SizeBytes: ConvertBytesToGb(disk.SizeBytes)})
				}

				validStorage := StorageDetails{
					Count: len(disks),
					Disk:  disks,
				}
				hardwareDetails[DISKLabel] = validStorage
			}

			if extractedProfile.NIC != nil {
				// Get the NIC details from the baremetal host and validate it into new structure
				var validNICS NICDetails

				for _, NIC := range host.Status.HardwareDetails.NIC {
					if NIC.PXE && CheckValidIP(NIC.IP) {
						validNICS.Nic.Name = NIC.Name
						validNICS.Nic.PXE = NIC.PXE
					}
				}

				validNICS.Count = len(host.Status.HardwareDetails.NIC)
				hardwareDetails[NICLabel] = validNICS
			}

			if extractedProfile.RAM != nil {
				// Get the RAM details from the baremetal host and validate it into new structure
				var RAM int64
				RAM = int64(host.Status.HardwareDetails.RAMMebibytes / 1024)
				hardwareDetails[RAMLabel] = RAM
			}

			if len(hardwareDetails) != 0 {
				validatedHostMap[host.ObjectMeta.Name] = hardwareDetails
			}
		}
	}
	return validatedHostMap
}
