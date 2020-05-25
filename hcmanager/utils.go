package hcmanager

import (
	"context"
	"errors"
	"net"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// DeleteLabels delete existing label of the baremetal host
func (mgr HardwareClassificationManager) DeleteLabels(ctx context.Context, hcMetaData v1.ObjectMeta, BMHList bmh.BareMetalHostList) error {

	// label for the baremetal host
	labelKey := LabelName + hcMetaData.Name

	// Delete existing labels for the same profile.
	for _, host := range BMHList.Items {

		if host.Status.Provisioning.State == Status {

			existingLabels := host.GetLabels()

			for key := range existingLabels {
				if key == labelKey {
					delete(existingLabels, key)
				}
			}

			host.SetLabels(existingLabels)

			err := mgr.client.Update(ctx, &host)
			if err != nil {
				return errors.New("Label Delete Failed")
			}
		}
	}

	return nil
}

//SetLabel update the labels of the comapred baremetal host
func (mgr HardwareClassificationManager) SetLabel(ctx context.Context, hcMetaData v1.ObjectMeta, comparedHost []string, BMHList bmh.BareMetalHostList, extractedLabels map[string]string) (bool, error, error) {

	// label for the baremetal host
	labelKey := LabelName + hcMetaData.Name
	setLabel := false

	for _, validHost := range comparedHost {
		for _, host := range BMHList.Items {
			m := make(map[string]string)
			if validHost == host.Name {
				// Getting all the existing labels on the matched host.
				availableLabels := host.GetLabels()
				mgr.Log.Info("Existing Labels ", validHost, availableLabels)

				for key, value := range availableLabels {
					m[key] = value
				}
				if extractedLabels != nil {
					for _, value := range extractedLabels {
						if value == "" {
							m[labelKey] = DefaultLabel
						} else {
							m[labelKey] = value
						}
					}
				} else {
					m[labelKey] = DefaultLabel
				}
				mgr.Log.Info("Labels to be applied ", validHost, m)

				// Setting labels on the matched host.
				host.SetLabels(m)
				err := mgr.client.Update(ctx, &host)
				if err != nil {
					return setLabel, errors.New("Failed to set label on host" + validHost + "Error :" + err.Error()), nil
				}
				setLabel = true
			}
		}
	}

	// delete existing labels
	if !setLabel {
		err := mgr.DeleteLabels(ctx, hcMetaData, BMHList)
		if err != nil {
			return setLabel, nil, errors.New(err.Error())
		}
	}

	return setLabel, nil, nil
}
