package hcmanager

import (
	"context"
	"errors"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

				return errors.New(host.Name + " Label Delete Failed")
			}
		}
	}

	return nil
}

//SetLabel update the labels of the comapred baremetal host
func (mgr HardwareClassificationManager) SetLabel(ctx context.Context, hcMetaData v1.ObjectMeta, comparedHost []string, BMHList bmh.BareMetalHostList, extractedLabels map[string]string) (bool, []string, error) {

	// label for the baremetal host
	labelKey := LabelName + hcMetaData.Name
	setLabel := false
	var errHost []string

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
					errHost = append(errHost, validHost+" "+err.Error())
				} else {
					setLabel = true
				}
			}
		}
	}

	// delete existing labels
	if !setLabel {
		err := mgr.DeleteLabels(ctx, hcMetaData, BMHList)
		if err != nil {
			return setLabel, errHost, errors.New(err.Error())
		}
	}

	return setLabel, errHost, nil
}
