/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ContainerSpec struct {
	Image string `json:"image"`
	Port  int32  `json:"port"`
}
type ServiceSpec struct {
	//+optional
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType string `json:"serviceType,omitempty"`
	//+optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
}

// AadeeSpec defines the desired state of Aadee
type AadeeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	DeploymentName string `json:"deploymentName,omitempty"`
	// Replicas defines number of pods will be running in the deployment
	Replicas  *int32        `json:"replicas"`
	Container ContainerSpec `json:"container"`
	// Service contains ServiceName, ServiceType, ServiceNodePort
	//+optional
	Service ServiceSpec `json:"service,omitempty"`
	// Foo is an example field of Aadee. Edit aadee_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
}

// AadeeStatus defines the observed state of Aadee
type AadeeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	AvailableReplicas *int32 `json:"availableReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Aadee is the Schema for the aadees API
type Aadee struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AadeeSpec   `json:"spec,omitempty"`
	Status AadeeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AadeeList contains a list of Aadee
type AadeeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aadee `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Aadee{}, &AadeeList{})
}
