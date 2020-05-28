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
	"strings"

	"github.com/go-logr/logr"

	hwcc "hardware-classification-controller/api/v1alpha1"
	utils "hardware-classification-controller/hcmanager"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	//HWControllerName Name to show in the logs
	HWControllerName = "HardwareClassification-Controller"
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

	// Initialize the logger with namespace
	//hcReconciler.Log = hcReconciler.Log.WithName(HWControllerName).WithValues("metal3-harwdwareclassification", req.NamespacedName)

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
	hcReconciler.Log.Info("Extracted hardware configurations successfully", "Profile", extractedProfile)

	// Get the new hardware classification manager
	hcManager := utils.NewHardwareClassificationManager(hcReconciler.Client, hcReconciler.Log)

	ErrValidation := hcManager.ValidateExtractedHardwareProfile(extractedProfile)

	if ErrValidation != nil {
		hcReconciler.Log.Error(ErrValidation, ErrValidation.Error())
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.ProfileMisConfigured, ErrValidation.Error(), hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	//Fetch baremetal host list for the given namespace
	hostList, BMHList, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		errMessage := "Unable to fetch BMH list from BMO"
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.FetchBMHListFailure, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	if len(hostList) == 0 {
		errMessage := "No BareMetalHost found in ready state"
		hcReconciler.Log.Info("No BareMetalHost found in ready state")
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.NoBMHHost, errMessage, hwcc.ProfileMatchStatusEmpty)
		return ctrl.Result{}, nil
	}

	//Extract the hardware details from the baremetal host list
	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)

	hcReconciler.Log.Info("Validated Hardware Details From Baremetal Hosts", "Validated Host List", validatedHardwareDetails)

	//Compare the host list with extracted profile and fetch the valid host names
	comparedHost := hcManager.MinMaxComparison(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	hcReconciler.Log.Info("Comapred Baremetal Hosts list Against User Profile ", "Compared Host Names", comparedHost)

	// set labels to matched hosts
	setLabel, errHost, deleteLabelErr := hcManager.SetLabel(ctx, hardwareClassification.ObjectMeta, comparedHost, BMHList, hardwareClassification.ObjectMeta.Labels)

	if len(errHost) > 0 {
		hcReconciler.handleErrorConditions(req, hardwareClassification, hwcc.LabelUpdateFailure, strings.Join(errHost, ","), hwcc.ProfileMatchStatusEmpty)
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
