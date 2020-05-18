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

	"github.com/pkg/errors"

	"github.com/go-logr/logr"

	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// HardwareClassificationReconciler reconciles a HardwareClassification object
type HardwareClassificationReconciler struct {
	Client client.Client
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

	// fetch BMH list from BMO
	validHostList, err := fetchBmhHostList(ctx, r, hardwareClassification.ObjectMeta.Namespace)

	if err != nil {
		//hardwareClassification.Status.ErrorType = hwcc.FetchBMHListFailure
		//hardwareClassification.Status.ErrorMessage = "Unable to fetch BMH list from BMO"
		//hardwareClassification.SetErrorMessage(hwcc.FetchBMHListFailure, "Unable to fetch BMH list from BMO")
		//r.Client.Status().Update(ctx, hardwareClassification)
		r.setErrorCondition(req, hardwareClassification, hwcc.FetchBMHListFailure, "Unable to fetch BMH list from BMO")
		//r.saveHWCCStatus(hardwareClassification)
		fmt.Println("Status Updated**********************", hardwareClassification.Status)
		return ctrl.Result{}, nil
	}

	if len(validHostList) == 0 {
		r.Log.Info("No BareMetal Host found in ready state")
		return ctrl.Result{}, nil
	}

	// Extract introspection data for each configuration provided in profile
	extractedHardwareDetails := extractHardwareDetails(extractedProfile, validHostList)

	r.Log.Info("Extracted hardware introspection details successfully", "IntrospectionDetails", extractedHardwareDetails)

	return ctrl.Result{}, nil
}

// fetchBmhHostList this function will fetch and return baremetal hosts in ready state
func fetchBmhHostList(ctx context.Context, r *HardwareClassificationReconciler, namespace string) ([]bmh.BareMetalHost, error) {

	bmhHostList := bmh.BareMetalHostList{}
	var validHostList []bmh.BareMetalHost

	opts := &client.ListOptions{
		Namespace: namespace,
	}

	// Get list of BareMetalHost from BMO
	err := r.Client.List(ctx, &bmhHostList, opts)
	err = errors.New("Unable to fetch BMH list")
	if err != nil {
		return validHostList, err
	}

	// Get hosts in ready status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
		}
	}

	return validHostList, nil
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

func (r *HardwareClassificationReconciler) setErrorCondition(request ctrl.Request, hardwareClassification *hwcc.HardwareClassification, errType hwcc.ErrorType, message string) (changed bool, err error) {
	var log = logf.Log.WithName("controller.HardwareClassification")

	reqLogger := log.WithValues("Request.Namespace",
		request.Namespace, "Request.Name", request.Name)

	changed = hardwareClassification.SetErrorMessage(errType, message)
	if changed {
		reqLogger.Info(
			"adding error message",
			"message", message,
		)
		err = r.saveHWCCStatus(hardwareClassification)
		if err != nil {
			err = errors.Wrap(err, "failed to update error message")
		}
	}

	return
}

func (r *HardwareClassificationReconciler) saveHWCCStatus(hcc *hwcc.HardwareClassification) error {
	//t := metav1.Now()
	//host.Status.LastUpdated = &t

	/*if err := r.saveHostAnnotation(host); err != nil {
		return err
	}*/

	//Refetch hwcc again
	obj := hcc.Status.DeepCopy()
	err := r.Client.Get(context.TODO(),

		client.ObjectKey{
			Name:      hcc.Name,
			Namespace: hcc.Namespace,
		},
		hcc,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to update Status")
	}

	hcc.Status = *obj
	err = r.Client.Status().Update(context.TODO(), hcc)
	fmt.Println("HCC after update*************", hcc.Status)
	return err
}

// SetupWithManager will add watches for this controller
func (r *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			handler.Funcs{},
		).
		Complete(r)
}
