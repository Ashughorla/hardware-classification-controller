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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HardwareClassificationSpec defines the desired state of HardwareClassification
type HardwareClassificationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ExpectedHardwareConfiguration defines expected hardware configurations for CPU, RAM, Disk, NIC.
	ExpectedHardwareConfiguration ExpectedHardwareConfiguration `json:"expectedValidationConfiguration"`
}

// ExpectedHardwareConfiguration details to match with the host
type ExpectedHardwareConfiguration struct {
	// +optional
	CPU *CPU `json:"CPU"`
	// +optional
	Disk *Disk `json:"Disk"`
	// +optional
	NIC *NIC `json:"NIC"`
	// +optional
	RAM *RAM `json:"RAM"`
}

// CPU contains CPU details extracted from the hardware profile
type CPU struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum cpu count, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum cpu count, greater than 0"`
	// +optional
	// +kubebuilder:validation:Pattern=`^(0\.\d*[1-9]\d*|[1-9]\d*(\.\d+)?)$`
	MinimumSpeed string `json:"minimumSpeed" description:"minimum speed of cpu, greater than 0"`
	// +optional
	// +kubebuilder:validation:Pattern=`^(0\.\d*[1-9]\d*|[1-9]\d*(\.\d+)?)$`
	MaximumSpeed string `json:"maximumSpeed" description:"maximum speed of cpu, greater than 0"`
}

// Disk contains disk details extracted from the hardware profile
type Disk struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum count of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumIndividualSizeGB int64 `json:"minimumIndividualSizeGB" description:"minimum individual size of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum count of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumIndividualSizeGB int64 `json:"maximumIndividualSizeGB" description:"maximum individual size of disk, greater than 0"`
}

// NIC contains nic details extracted from the hardware profile
type NIC struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum count of nics, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum count of nics, greater than 0"`
}

// RAM contains ram details extracted from the hardware profile
type RAM struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumSizeGB int `json:"minimumSizeGB" description:"minimun size of ram, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumSizeGB int `json:"maximumSizeGB" description:"maximum size of ram, greater than 0"`
}

// HardwareClassificationStatus defines the observed state of HardwareClassification
type HardwareClassificationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// HardwareClassification is the Schema for the hardwareclassifications API
type HardwareClassification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareClassificationSpec   `json:"spec,omitempty"`
	Status HardwareClassificationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HardwareClassificationList contains a list of HardwareClassification
type HardwareClassificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HardwareClassification `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HardwareClassification{}, &HardwareClassificationList{})
}
