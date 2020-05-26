package hcmanager

import (
	"context"
	"errors"
	"fmt"
	"net"

	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	//LabelName initial name to the label key as hardware classification group
	LabelName = "hardwareclassification.metal3.io/"

	//Status extract the baremetal host for status ready
	Status = "ready"

	//DefaultLabel if label is missing from the Extracted Hardware Profile
	DefaultLabel = "matches"

	//CPULable label for extraction of hardware details
	CPULable = "CPU"

	//NICLable label for extraction of hardware details
	NICLable = "NIC"

	//DISKLable label for extraction of hardware details
	DISKLable = "DISK"

	//RAMLable label for extraction of hardware details
	RAMLable = "RAM"
)

//FetchBmhHostList this function will fetch and return baremetal hosts in ready state
func (mgr HardwareClassificationManager) FetchBmhHostList(Namespace string) ([]bmh.BareMetalHost, bmh.BareMetalHostList, error) {

	ctx := context.Background()

	bmhHostList := bmh.BareMetalHostList{}
	var validHostList []bmh.BareMetalHost

	opts := &client.ListOptions{
		Namespace: Namespace,
	}

	// Get list of BareMetalHost from BMO
	err := mgr.client.List(ctx, &bmhHostList, opts)
	if err != nil {
		return validHostList, bmhHostList, errors.New(err.Error())
	}

	// Get hosts in ready status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
		}
	}

	return validHostList, bmhHostList, nil
}

//CheckValidIP uses net package to check if the IP is valid or not
func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

//ConvertBytesToGb it converts the Byte into GB
func ConvertBytesToGb(inBytes int64) int64 {
	inGb := (inBytes / 1024 / 1024 / 1024)
	return inGb
}

//ValidateExtractedHardwareProfile it will validate the extracted hardware profile and log the warnings
func (mgr HardwareClassificationManager) ValidateExtractedHardwareProfile(extractedProfile hwcc.ExpectedHardwareConfiguration) error {

	if (extractedProfile.CPU == nil) &&
		(extractedProfile.RAM == nil) &&
		(extractedProfile.Disk == nil) &&
		(extractedProfile.NIC == nil) {
		return errors.New("Expected Profile details can not be empty")
	}

	fmt.Print("CPU :", extractedProfile.CPU)
	fmt.Print("RAM :", extractedProfile.RAM)
	fmt.Print("DISK :", extractedProfile.Disk)
	fmt.Print("NIC :", extractedProfile.NIC)

	if extractedProfile.CPU == nil {
		mgr.Log.Info("WARNING CPU details are empty")
	} else {
		if (extractedProfile.CPU.MaximumCount > 0) ||
			(extractedProfile.CPU.MinimumCount > 0) ||
			(extractedProfile.CPU.MaximumSpeed != "") ||
			(extractedProfile.CPU.MinimumSpeed != "") {
		} else {
			return errors.New("Extracted CPU details can not be empty")
		}
	}

	if extractedProfile.RAM == nil {
		mgr.Log.Info("WARNING RAM details are empty")
	} else {
		if (extractedProfile.RAM.MaximumSizeGB == 0) ||
			(extractedProfile.RAM.MinimumSizeGB == 0) {
			return errors.New("Extracted RAM details can not be empty")
		}

	}

	if extractedProfile.Disk == nil {
		mgr.Log.Info("WARNING DISK details are empty")
	} else {
		if (extractedProfile.Disk.MaximumCount == 0) ||
			(extractedProfile.Disk.MinimumCount == 0) ||
			(extractedProfile.Disk.MaximumIndividualSizeGB == 0) ||
			(extractedProfile.Disk.MinimumIndividualSizeGB == 0) {
			return errors.New("Extracted DISK details can not be empty")
		}

	}

	if extractedProfile.NIC == nil || (extractedProfile.NIC == &hwcc.NIC{}) {
		mgr.Log.Info("WARNING NIC details is empty")
	} else {

		if (extractedProfile.NIC.MaximumCount == 0) ||
			(extractedProfile.NIC.MinimumCount == 0) {
			return errors.New("Extracted DISK details can not be empty")
		}
	}

	return nil
}
