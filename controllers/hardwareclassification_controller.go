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

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	hwcc "hardware-classification-controller/api/v1alpha1"
	utils "hardware-classification-controller/hcmanager"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
func (hcReconciler *HardwareClassificationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassification{}

	if err := hcReconciler.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration

	if (extractedProfile.CPU == &hwcc.CPU{}) && (extractedProfile.RAM == &hwcc.RAM{}) && (extractedProfile.Disk == &hwcc.Disk{}) && (extractedProfile.NIC == &hwcc.NIC{}) {
		hcReconciler.Log.Info("Expected Profile details can not be empty")
		errMessage := "Expected Profile details can not be empty"
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.ProfileMisConfigured, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	hcReconciler.Log.Info("Extracted hardware configurations successfully", "Profile", extractedProfile)

	// Get the new hardware classification manager
	hcManager := utils.NewHardwareClassificationManager(hcReconciler.Client, hcReconciler.Log)

	//Fetch baremetal host list for the given namespace
	hostList, BMHList, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		errMessage := "Unable to fetch BMH list from BMO"
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.FetchBMHListFailure, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	if len(hostList) == 0 {
		hcReconciler.Log.Info("No BareMetalHost found in ready state")
		return ctrl.Result{}, nil
	}

	//Extract the hardware details from the baremetal host list
	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)

	hcReconciler.Log.Info("Validated Hardware Details From Baremetal Hosts", "Validated Host List", validatedHardwareDetails)

	//Comapre the host list with extracted profile and fetch the valid host names
	comparedHost := hcManager.MinMaxComparison(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	hcReconciler.Log.Info("Comapred Baremetal Hosts list Against User Profile ", "Compared Host Names", comparedHost)

	// set labels to matched hosts
	setLabel, setLabelErr, deleteLabelErr := hcManager.SetLabel(ctx, hardwareClassification.ObjectMeta, comparedHost, BMHList, hardwareClassification.ObjectMeta.Labels)

	if setLabelErr != nil {
		errMessage := "Failed to set labels on BareMetalHost"
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.LabelUpdateFailure, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	if deleteLabelErr != nil {
		errMessage := "Failed to delete existing labels on BareMetalHost"
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.LabelDeleteFailure, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	if setLabel {
		hcReconciler.updateProfileMatchStatus(req, hardwareClassification, hwcc.ProfileMatchStatusMatched)
	} else {
		hcReconciler.updateProfileMatchStatus(req, hardwareClassification, hwcc.ProfileMatchStatusUnMatched)
	}

	return ctrl.Result{}, nil
}

func (hcReconciler *HardwareClassificationReconciler) updateProfileMatchStatus(req ctrl.Request, hc *hwcc.HardwareClassification, status hwcc.ProfileMatchStatus) {
	hcReconciler.Log.Info("Current status is:", "ProfileMatchStatus", status)
	if hc.SetProfileMatchStatus(status) {
		hcReconciler.Log.Info("clearing previous error message")
		hc.ClearError()

		hcReconciler.Log.Info("Upating status as:", "ProfileMatchStatus", status)
		err := hcReconciler.saveHWCCStatus(hc)
		if err != nil {
			hcReconciler.Log.Error(err, "Error while saving ProfileMatchStatus")
		}
	}
}

func (hcReconciler *HardwareClassificationReconciler) handleErrorConditions(req ctrl.Request, hc *hwcc.HardwareClassification, errorType hwcc.ErrorType, message string, status hwcc.ProfileMatchStatus) {
	hcReconciler.setErrorCondition(req, hc, hwcc.FetchBMHListFailure, message)

	hcReconciler.Log.Info("Upating status as:", "ProfileMatchStatus", status)
	hc.SetProfileMatchStatus(status)
	err := hcReconciler.saveHWCCStatus(hc)
	if err != nil {
		hcReconciler.Log.Error(err, "Error while saving ProfileMatchStatus")
	}
}

func (hcReconciler *HardwareClassificationReconciler) setErrorCondition(request ctrl.Request, hardwareClassification *hwcc.HardwareClassification, errType hwcc.ErrorType, message string) (changed bool, err error) {
	var log = logf.Log.WithName("controller.HardwareClassification")

	reqLogger := log.WithValues("Request.Namespace",
		request.Namespace, "Request.Name", request.Name)

	changed = hardwareClassification.SetErrorMessage(errType, message)
	if changed {
		reqLogger.Info(
			"adding error message",
			"message", message,
		)
		err = hcReconciler.saveHWCCStatus(hardwareClassification)
		if err != nil {
			err = errors.Wrap(err, "failed to update error message")
		}
	} else {
		reqLogger.Info(
			"aleady added error message",
			"message", message,
		)
	}

	return
}

func (hcReconciler *HardwareClassificationReconciler) saveHWCCStatus(hcc *hwcc.HardwareClassification) error {

	//Refetch hwcc again
	obj := hcc.Status.DeepCopy()
	err := hcReconciler.Client.Get(context.TODO(),

		client.ObjectKey{
			Name:      hcc.Name,
			Namespace: hcc.Namespace,
		},
		hcc,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to update HardwareClassification Status")
	}

	hcc.Status = *obj
	err = hcReconciler.Client.Status().Update(context.TODO(), hcc)
	return err
}

// SetupWithManager will add watches for this controller
func (hcReconciler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			handler.Funcs{},
		).
		Complete(hcReconciler)
}
