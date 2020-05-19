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
func (hcReconiler *HardwareClassificationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassification{}

	if err := hcReconiler.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration
	hcReconiler.Log.Info("Extracted hardware configurations successfully", "Profile", extractedProfile)

	//Get the new hardware classification manager
	hcManager := utils.NewHardwareClassificationManager(hcReconiler.Client, hcReconiler.Log)

	//Fetch baremetal host list for the given namespace
	hostList, BMHList, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		hcReconiler.setErrorCondition(req, hardwareClassification, hwcc.FetchBMHListFailure, "Unable to fetch BMH list from BMO")
		return ctrl.Result{}, nil
	}

	if len(hostList) == 0 {
		hcReconiler.Log.Info("No BareMetalHost found in ready state")
		return ctrl.Result{}, nil
	}

	//Extract the hardware details from the baremetal host list
	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)
	hcReconiler.Log.Info("Validated Hardware Details From Baremetal Hosts", "Validated Host List", validatedHardwareDetails)

	//Comapre the host list with extracted profile and fetch the valid host names
	comparedHost := hcManager.MinMaxComparison(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	hcReconiler.Log.Info("Comapred Baremetal Hosts list Against User Profile ", "Compared Host Names", comparedHost)

	//Delete the existing label if any to add new label
	err = hcManager.DeleteLabels(ctx, hardwareClassification.ObjectMeta, BMHList)
	if err != nil {
		hcReconiler.setErrorCondition(req, hardwareClassification, hwcc.LabelDeleteFailure, "Failed to delete existing labels of BareMetalHost")
		//hcReconiler.Log.Error(err, "Deleting Existing Baremetal Host Label Failed")
		return ctrl.Result{}, nil
	}

	//set new labels to the valid host
	hcManager.SetLabel(ctx, hardwareClassification.ObjectMeta, comparedHost, BMHList, hardwareClassification.ObjectMeta.Labels)
	if err != nil {
		hcReconiler.setErrorCondition(req, hardwareClassification, hwcc.LabelUpdateFailure, "Updating Baremetal Host Label Failed")
		//hcReconiler.Log.Error(err, "Updating Baremetal Host Label Failed")
		return ctrl.Result{}, nil
	}

	// Reaching this point means the profile matched to one of BaremetalHost,
	// so clear any previous error and record the matched status in the
	// status block.
	hcReconiler.Log.Info(Info("updating ProfileMatchStatus status field")
	hcReconiler.Log.Info("clearing previous error message")

	hardwareClassification.ClearError()
	hardwareClassification.Status.ProfileMatchStatus = "matched"
	r.saveHWCCStatus(hardwareClassification)

	return ctrl.Result{}, nil
}

func (hcReconiler *HardwareClassificationReconciler) setErrorCondition(request ctrl.Request, hardwareClassification *hwcc.HardwareClassification, errType hwcc.ErrorType, message string) (changed bool, err error) {
	var log = logf.Log.WithName("controller.HardwareClassification")

	reqLogger := log.WithValues("Request.Namespace",
		request.Namespace, "Request.Name", request.Name)

	changed = hardwareClassification.SetErrorMessage(errType, message)
	if changed {
		reqLogger.Info(
			"adding error message",
			"message", message,
		)
		err = hcReconiler.saveHWCCStatus(hardwareClassification)
		if err != nil {
			err = errors.Wrap(err, "failed to update error message")
		}
	}

	return
}

func (hcReconiler *HardwareClassificationReconciler) saveHWCCStatus(hcc *hwcc.HardwareClassification) error {
	//t := metav1.Now()
	//host.Status.LastUpdated = &t

	//Refetch hwcc again
	obj := hcc.Status.DeepCopy()
	err := hcReconiler.Client.Get(context.TODO(),

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
	err = hcReconiler.Client.Status().Update(context.TODO(), hcc)
	return err
}

// SetupWithManager will add watches for this controller
func (hcReconiler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			handler.Funcs{},
		).
		Complete(hcReconiler)
}
