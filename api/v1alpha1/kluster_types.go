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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ContainerSpec struct {
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

type ServiceSpec struct {
	//+optional
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType string `json:"serviceType"`
	//+optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
	ServicePort     int32 `json:"servicePort"`
}

// PritamSpec defines the desired state of Pritam
type PritamSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Name      string        `json:"name,omitempty"`
	Replicas  *int32        `json:"replicas"`
	Container ContainerSpec `json:"container"`
	Service   ServiceSpec   `json:"service,omitempty"`
}

// PritamStatus defines the observed state of Pritam
type PritamStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvilableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pritam is the Schema for the Pritams API
type Pritam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PritamSpec   `json:"spec,omitempty"`
	Status PritamStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PritamList contains a list of Pritam
type PritamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pritam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pritam{}, &PritamList{})
}
