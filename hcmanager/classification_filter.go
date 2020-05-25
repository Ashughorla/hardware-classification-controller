package hcmanager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"
	"strconv"
)

// MinMaxComparison it will compare the minimum and maximum comparison based on the value provided by the user and check for the valid host
func (mgr HardwareClassificationManager) MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) []string {

	var comparedHost []string

	for hostname, details := range validatedHost {
		isHostValid := true

		for _, value := range details {

			isValid := false

			cpu, CPUOK := value.(CPU)
			if CPUOK {
				if (expectedHardwareprofile.CPU.MaximumCount > 0) ||
					(expectedHardwareprofile.CPU.MinimumCount > 0) ||
					(expectedHardwareprofile.CPU.MaximumSpeed != "") ||
					(expectedHardwareprofile.CPU.MinimumSpeed != "") {
					if checkCPUCount(mgr, cpu, expectedHardwareprofile.CPU) {
						isValid = true
					}
				}

			}

			ram, RAMOK := value.(RAM)
			if RAMOK {
				if (expectedHardwareprofile.RAM.MaximumSizeGB > 0) ||
					(expectedHardwareprofile.RAM.MinimumSizeGB > 0) {
					if checkRAM(mgr, ram, expectedHardwareprofile.RAM) {
						isValid = true
					}
				}

			}

			nics, NICSOK := value.(NIC)
			if NICSOK {
				if (expectedHardwareprofile.NIC.MaximumCount > 0) ||
					(expectedHardwareprofile.NIC.MinimumCount > 0) {
					if checkNICS(mgr, nics, expectedHardwareprofile.NIC) {
						isValid = true
					}
				}
			}

			disk, DISKOK := value.(Storage)
			if DISKOK {
				if (expectedHardwareprofile.Disk.MaximumCount > 0) ||
					(expectedHardwareprofile.Disk.MinimumCount > 0) ||
					(expectedHardwareprofile.Disk.MaximumIndividualSizeGB > 0) ||
					(expectedHardwareprofile.Disk.MinimumIndividualSizeGB > 0) {
					if checkDiskDetails(mgr, disk, expectedHardwareprofile.Disk) {
						isValid = true
					}
				}

			}

			if !isValid {
				isHostValid = false
				break
			}

		}

		if isHostValid {
			comparedHost = append(comparedHost, hostname)
			mgr.Log.Info(hostname, " Matches profile ", ProfileName)

		} else {
			mgr.Log.Info(hostname, " Did not matches profile ", ProfileName)
		}

	}

	return comparedHost

}

//checkCPUCount this function checks the CPU details for both min and max parameters
func checkCPUCount(mgr HardwareClassificationManager, cpu CPU, expectedCPU *hwcc.CPU) bool {

	if (expectedCPU.MaximumCount > 0) && (expectedCPU.MinimumCount > 0) {

		mgr.Log.Info("", "Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		mgr.Log.Info("", "Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if (expectedCPU.MinimumCount > cpu.Count) || (expectedCPU.MaximumCount < cpu.Count) {
			mgr.Log.Info("CPU MINMAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MaximumCount > 0 {

		mgr.Log.Info("", "Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MaximumCount < cpu.Count {

			mgr.Log.Info("CPU MAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MinimumCount > 0 {

		mgr.Log.Info("", "Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MinimumCount > cpu.Count {

			mgr.Log.Info("CPU MIN COUNT did not match")
			return false
		}

	}

	if (expectedCPU.MaximumSpeed != "") && (expectedCPU.MinimumSpeed != "") {
		MaxSpeed, errMax := strconv.ParseFloat(expectedCPU.MaximumSpeed, 64)
		MinSpeed, errMin := strconv.ParseFloat(expectedCPU.MinimumSpeed, 64)
		if errMax == nil && errMin == nil {

			mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MinSpeed > 0 && MaxSpeed > 0 {
				if (MinSpeed > cpu.ClockSpeed) || (MaxSpeed < cpu.ClockSpeed) {

					mgr.Log.Info("CPU MINMAX ClockSpeed did not match")
					return false
				}

			}
		}

	} else if expectedCPU.MaximumSpeed != "" {
		MaxSpeed, errMax := strconv.ParseFloat(expectedCPU.MaximumSpeed, 64)
		if errMax == nil {

			mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MaxSpeed > 0 {
				if MaxSpeed < cpu.ClockSpeed {

					mgr.Log.Info("CPU MAX ClockSpeed did not match")
					return false
				}

			}
		}
	} else if expectedCPU.MinimumSpeed != "" {
		MinSpeed, errMin := strconv.ParseFloat(expectedCPU.MinimumSpeed, 64)
		if errMin == nil {

			mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MinSpeed > 0 {
				if MinSpeed > cpu.ClockSpeed {

					mgr.Log.Info("CPU MIN ClockSpeed did not match")
					return false
				}

			}
		}
	}

	return true

}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(mgr HardwareClassificationManager, nics NIC, expectedNIC *hwcc.NIC) bool {

	if (expectedNIC.MaximumCount > 0) && (expectedNIC.MinimumCount > 0) {

		mgr.Log.Info("", "Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		mgr.Log.Info("", "Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if (expectedNIC.MinimumCount > nics.Count) || (expectedNIC.MaximumCount < nics.Count) {

			mgr.Log.Info("NICS MINMAX count did not match")
			return false
		}
	} else if expectedNIC.MaximumCount > 0 {

		mgr.Log.Info("", "Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if expectedNIC.MaximumCount < nics.Count {

			mgr.Log.Info("NICS MAX count did not match")
			return false
		}

	} else if expectedNIC.MinimumCount > 0 {

		mgr.Log.Info("", "Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		if expectedNIC.MinimumCount > nics.Count {

			mgr.Log.Info("NICS MIN count did not match")
			return false
		}

	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(mgr HardwareClassificationManager, ram RAM, expectedRAM *hwcc.RAM) bool {
	if (expectedRAM.MaximumSizeGB > 0) && (expectedRAM.MinimumSizeGB > 0) {

		mgr.Log.Info("", "Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		mgr.Log.Info("", "Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if (expectedRAM.MinimumSizeGB > ram.RAMGb) || (expectedRAM.MaximumSizeGB < ram.RAMGb) {

			mgr.Log.Info("RAM MINMAX SIZE did not match")
			return false
		}
	} else if expectedRAM.MaximumSizeGB > 0 {

		mgr.Log.Info("", "Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MaximumSizeGB < ram.RAMGb {

			mgr.Log.Info("RAM MAX SIZE did not match")
			return false
		}

	} else if expectedRAM.MinimumSizeGB > 0 {

		mgr.Log.Info("", "Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MinimumSizeGB > ram.RAMGb {

			mgr.Log.Info("RAM MIN SIZE did not match")
			return false
		}

	}
	return true
}

//checkDiskDetails this function checks the Disk details for both min and max parameters
func checkDiskDetails(mgr HardwareClassificationManager, storage Storage, expectedDisk *hwcc.Disk) bool {

	if (expectedDisk.MaximumCount > 0) && (expectedDisk.MinimumCount > 0) {
		mgr.Log.Info("", "Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", storage.Count)
		mgr.Log.Info("", "Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", storage.Count)

		if (expectedDisk.MinimumCount > storage.Count) || (expectedDisk.MaximumCount < storage.Count) {
			mgr.Log.Info("Disk MINMAX Count did not match")
			return false
		}

	} else if expectedDisk.MaximumCount > 0 {
		if expectedDisk.MaximumCount < storage.Count {
			mgr.Log.Info("Disk MAX Count did not match")
			return false
		}
	} else if expectedDisk.MinimumCount > 0 {
		if expectedDisk.MinimumCount > storage.Count {
			mgr.Log.Info("Disk MIN Count did not match")
			return false
		}

	}

	for _, disk := range storage.Disk {
		if expectedDisk.MaximumIndividualSizeGB > 0 && expectedDisk.MinimumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			mgr.Log.Info("", "Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			if (expectedDisk.MaximumIndividualSizeGB < disk.SizeGb) || (expectedDisk.MinimumIndividualSizeGB > disk.SizeGb) {

				mgr.Log.Info("Disk MINMAX SIZE did not match")
				return false
			}
		} else if expectedDisk.MaximumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {

				mgr.Log.Info("Disk MAX SIZE did not match")
				return false
			}
		} else if expectedDisk.MinimumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {

				mgr.Log.Info("Disk MIN SIZE did not match")
				return false
			}
		}
	}

	return true
}
