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
	utils "hardware-classification-controller/hcutils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

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

	hcManager := utils.NewHardwareClassificationManager(hcReconiler.Client, hcReconiler.Log)

	hostList, _, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		hcReconiler.Log.Error(err, "Fetch Baremetal Host List Failed")
		return ctrl.Result{}, nil
	}

	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)
	hcReconiler.Log.Info("Validated Hardware Details From Baremetal Hosts", "Validated Host List", validatedHardwareDetails)

	// comparedHost := hcManager.MinMaxComparison(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	// hcReconiler.Log.Info("Comapred Baremetal Hosts list Against User Profile ", comparedHost)

	// err = hcManager.DeleteLabels(ctx, hardwareClassification.ObjectMeta, BMHList)
	// if err != nil {
	// 	hcReconiler.Log.Error(err, "Deleting Existing Baremetal Host Label Failed")
	// 	return ctrl.Result{}, nil
	// }

	// hcManager.SetLabel(ctx, hardwareClassification.ObjectMeta, comparedHost, BMHList, hardwareClassification.ObjectMeta.Labels)
	// if err != nil {
	// 	hcReconiler.Log.Error(err, "Updating Baremetal Host Label Failed")
	// 	return ctrl.Result{}, nil
	// }

	hardwareClassification = &hwcc.HardwareClassification{}
	return ctrl.Result{}, nil
}

// SetupWithManager will add watches for this controller
func (hcReconiler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(hcReconiler.WatchHardwareClassification),
			},
		).
		Complete(hcReconiler)
}

// WatchHardwareClassification will return a reconcile request for a
// HardwareClassification if the event is for a HardwareClassification.
func (hcReconiler *HardwareClassificationReconciler) WatchHardwareClassification(obj handler.MapObject) []ctrl.Request {
	if profile, ok := obj.Object.(*hwcc.HardwareClassification); ok {
		fmt.Println("In Watcher Function for name: **************", profile.ObjectMeta.Name)
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
