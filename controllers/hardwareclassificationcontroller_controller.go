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

	hwcc "hardware-classification-controller/api/v1alpha1"

	validation "hardware-classification-controller/validate"

	"github.com/go-logr/logr"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HardwareClassificationControllerReconciler reconciles a HardwareClassificationController object
type HardwareClassificationControllerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconcile function
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers/status,verbs=get;update;patch
func (r *HardwareClassificationControllerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassificationController{}
	if err := r.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfileList := hardwareClassification.Spec.ExpectedHardwareConfiguration

	// Get a list of BaremetalHost from Baremetal-Operator and metal3 namespace
	bmhHostList, err := fetchBmhHostList(ctx, r, hardwareClassification.Spec.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Pass the baremetal host list and extractedhardwareProfile to comparison function
	validHostList := validation.Comparison(bmhHostList, extractedProfileList)

	fmt.Println("Hostname and valid profile details:-")
	fmt.Println("*********************************************")
	for hostDetails, profile := range validHostList {

		host, ok := hostDetails.(*bmh.HardwareDetails)
		if ok {
			fmt.Println("Host :- ", host.Hostname)
		}

		fmt.Println("Profiles")
		fmt.Println("")
		fmt.Printf("%+v \n\n", profile)
		fmt.Println("*********************************************")
	}

	return ctrl.Result{}, nil
}

// fetchBmhHostList this function will fetch the baremetal host list from the baremetal operator and return it to the reconciler functioni
func fetchBmhHostList(ctx context.Context, r *HardwareClassificationControllerReconciler, namespace string) ([]bmh.BareMetalHost, error) {

	bmhHostList := bmh.BareMetalHostList{}
	validHostList := []bmh.BareMetalHost{}
	hardwareClassification := &hwcc.HardwareClassificationController{}

	opts := &client.ListOptions{
		Namespace: namespace,
	}

	// Get list of BareMetalHost
	err := r.Client.List(ctx, &bmhHostList, opts)
	if err != nil {
		setError(hardwareClassification, "Failed to get BareMetalHost List")
		return nil, err
	}

	// Get hosts in ready and inspecting status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" || host.Status.Provisioning.State == "inspecting" {
			validHostList = append(validHostList, host)
		}

	}

	return validHostList, nil
}

// setError sets the ErrorMessage field on the baremetalmachine
func setError(hwcc *hwcc.HardwareClassificationController, message string) {
	hwcc.Status.ErrorMessage = pointer.StringPtr(message)
}

//SetupWithManager SetupWithManager
func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassificationController{}).
		Complete(r)
}
