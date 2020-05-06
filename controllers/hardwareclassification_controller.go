/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	hwcc "hardware-classification-controller/api/v1alpha1"
	filter "hardware-classification-controller/classification_filter"
	validate "hardware-classification-controller/validation"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// HardwareClassificationReconciler reconciles a HardwareClassification object
type HardwareClassificationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconcile function
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassifications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassifications/status,verbs=get;update;patch
func (r *HardwareClassificationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassification{}

	if err := r.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration
	r.Log.Info("Extracted hardware configurations successfully", "Profile", extractedProfile)
	extractedLabels := hardwareClassification.ObjectMeta.Labels

	hardwareClassification1 := hwcc.HardwareClassification{}
	// fetch BMH list from BMO
	validHostList := fetchBmhHostList(ctx, r, hardwareClassification1.ObjectMeta.Namespace)

	// Extract introspection data for each configuration provided in profile
	extractedHardwareDetails := extractHardwareDetails(extractedProfile, validHostList)

	r.Log.Info("Extracted hardware introspection details successfully", "IntrospectionDetails", extractedHardwareDetails)

	if len(extractedHardwareDetails) > 0 {
		validatedHardwareDetails := validate.Validation(extractedHardwareDetails)
		comparedHost := filter.MinMaxComparison(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
		fmt.Println("List of Comapred Host", comparedHost)
		r.Log.Info("Validated Host from comparison function", "Validated host", extractedHardwareDetails)
		setValidLabel(ctx, r, hardwareClassification.ObjectMeta.Name, comparedHost, extractedLabels)
	} else {
		fmt.Println("Provided configurations are not valid")
	}
	hardwareClassification = &hwcc.HardwareClassification{}
	return ctrl.Result{}, nil
}

// setValidLabel will add label to the hosts which matched ExpectedHardwareConfiguraton
func setValidLabel(ctx context.Context, r *HardwareClassificationReconciler, profileName string, matchedHosts []string, extractedLabels map[string]string) {
	// Get updated object to set labels on
	bmhHostList := bmh.BareMetalHostList{}
	opts := &client.ListOptions{
		Namespace: "metal3",
	}
	r.Log.Info("Getting updated host list to set labels")
	err := r.Client.List(ctx, &bmhHostList, opts)
	if err != nil {
		fmt.Println("Error while getting updated host list for labels : ", err)
	} else {
		r.Log.Info("Received updated host list to set labels")
	}

	labelKey := "hardwareclassification.metal3.io/" + profileName

	// Delete existing labels for the same profile.
	r.Log.Info("Checking if labels are already present for this profile")
	for i, _ := range bmhHostList.Items {
		existingLabels := bmhHostList.Items[i].GetLabels()
		for key, _ := range existingLabels {
			if key == labelKey {
				delete(existingLabels, key)
			}
		}
		bmhHostList.Items[i].SetLabels(existingLabels)
		err = r.Client.Update(ctx, &bmhHostList.Items[i])
		if err != nil {
			r.Log.Error(err, "Failed to delete existing labels with this profile")
		}
	}

	// Set latest labels on the matching hosts
	for _, validHost := range matchedHosts {
		for i, host := range bmhHostList.Items {
			m := make(map[string]string)
			if validHost == host.Name {
				// Getting all the existing labels on the matched host.
				availableLabels := bmhHostList.Items[i].GetLabels()
				r.Log.Info("Existing Labels ", validHost, availableLabels)
				for key, value := range availableLabels {
					m[key] = value
				}
				if extractedLabels != nil {
					for _, value := range extractedLabels {
						if value == "" {
							m[labelKey] = "matches"
						} else {
							m[labelKey] = value
						}
					}
				} else {
					m[labelKey] = "matches"
				}
				r.Log.Info("Labels to be applied ", validHost, m)

				// Setting labels on the matched host.
				bmhHostList.Items[i].SetLabels(m)
				err = r.Client.Update(ctx, &bmhHostList.Items[i])
				if err != nil {
					r.Log.Error(err, "Failed to set labels")
				} else {
					r.Log.Info("Labels updated successfully")
				}
			}
		}
	}
}

// fetchBmhHostList this function will fetch and return baremetal hosts in ready state
func fetchBmhHostList(ctx context.Context, r *HardwareClassificationReconciler, namespace string) []bmh.BareMetalHost {

	bmhHostList := bmh.BareMetalHostList{}
	validHostList := []bmh.BareMetalHost{}

	opts := &client.ListOptions{
		Namespace: namespace,
	}

	// Get list of BareMetalHost from BMO
	err := r.Client.List(ctx, &bmhHostList, opts)
	if err != nil {
		r.Log.Error(nil, "Unable to extract details", "error", err.Error())
		return validHostList
	}

	// Get hosts in ready status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
		}
	}

	return validHostList
}

// extractHardwareDetails this function will return map containing
// introspection details for a host.
func extractHardwareDetails(extractedProfile hwcc.ExpectedHardwareConfiguration,
	bmhList []bmh.BareMetalHost) map[string]map[string]interface{} {

	extractedHardwareDetails := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range bmhList {
			introspectionDetails := make(map[string]interface{})

			if extractedProfile.CPU != nil {
				introspectionDetails["CPU"] = host.Status.HardwareDetails.CPU
			}

			if extractedProfile.Disk != nil {
				introspectionDetails["Disk"] = host.Status.HardwareDetails.Storage
			}

			if extractedProfile.NIC != nil {
				introspectionDetails["NIC"] = host.Status.HardwareDetails.NIC
			}

			if extractedProfile.RAM != nil {
				introspectionDetails["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
			}

			if len(introspectionDetails) > 0 {
				extractedHardwareDetails[host.ObjectMeta.Name] = introspectionDetails
			}
		}
	}
	return extractedHardwareDetails
}

// SetupWithManager will add watches for this controller
func (r *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.WatchHardwareClassification),
			},
		).
		Complete(r)
}

// WatchHardwareClassification will return a reconcile request for a
// HardwareClassification if the event is for a HardwareClassification.
func (r *HardwareClassificationReconciler) WatchHardwareClassification(obj handler.MapObject) []ctrl.Request {
	if profile, ok := obj.Object.(*hwcc.HardwareClassification); ok {
		return []ctrl.Request{
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      profile.ObjectMeta.Name,
					Namespace: profile.ObjectMeta.Namespace,
				},
			},
		}
	}
	return []ctrl.Request{}
}
