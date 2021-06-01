/*
Copyright 2021.

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

// +kubebuilder:validation:Optional
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type LRPSpec struct {
	// +kubebuilder:validation:Required
	GUID        string `json:"GUID"`
	Version     string `json:"version"`
	ProcessType string `json:"processType"`
	AppName     string `json:"appName"`
	AppGUID     string `json:"appGUID"`
	OrgName     string `json:"orgName"`
	OrgGUID     string `json:"orgGUID"`
	SpaceName   string `json:"spaceName"`
	SpaceGUID   string `json:"spaceGUID"`
	// +kubebuilder:validation:Required
	Image           string            `json:"image"`
	Command         []string          `json:"command,omitempty"`
	Sidecars        []Sidecar         `json:"sidecars,omitempty"`
	PrivateRegistry *PrivateRegistry  `json:"privateRegistry,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	Health          Healthcheck       `json:"health"`
	Ports           []int32           `json:"ports,omitempty"`
	// +kubebuilder:default:=1
	Instances int   `json:"instances"`
	MemoryMB  int64 `json:"memoryMB"`
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Required
	DiskMB int64 `json:"diskMB"`
	// +kubebuilder:validation:Format:=uint8
	CPUWeight              uint8             `json:"cpuWeight"`
	VolumeMounts           []VolumeMount     `json:"volumeMounts,omitempty"`
	LastUpdated            string            `json:"lastUpdated"`
	UserDefinedAnnotations map[string]string `json:"userDefinedAnnotations,omitempty"`
}

type LRPStatus struct {
	Replicas int32 `json:"replicas"`
}

type Route struct {
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
}

type Sidecar struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Command  []string          `json:"command"`
	MemoryMB int64             `json:"memoryMB"`
	Env      map[string]string `json:"env,omitempty"`
}

type PrivateRegistry struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VolumeMount struct {
	MountPath string `json:"mountPath"`
	ClaimName string `json:"claimName"`
}

type Healthcheck struct {
	Type     string `json:"type"`
	Port     int32  `json:"port"`
	Endpoint string `json:"endpoint"`
	// +kubebuilder:validation:Format:=uint8
	TimeoutMs uint `json:"timeoutMs"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LRP is the Schema for the lrps API
type LRP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LRPSpec   `json:"spec,omitempty"`
	Status LRPStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LRPList contains a list of LRP
type LRPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LRP `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LRP{}, &LRPList{})
}
