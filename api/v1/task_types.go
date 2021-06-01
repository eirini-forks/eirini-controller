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

type TaskSpec struct {
	// +kubebuilder:validation:Required
	GUID string `json:"GUID"`
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Image           string            `json:"image"`
	PrivateRegistry *PrivateRegistry  `json:"privateRegistry,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	// +kubebuilder:validation:Required
	Command   []string `json:"command,omitempty"`
	AppName   string   `json:"appName"`
	AppGUID   string   `json:"appGUID"`
	OrgName   string   `json:"orgName"`
	OrgGUID   string   `json:"orgGUID"`
	SpaceName string   `json:"spaceName"`
	SpaceGUID string   `json:"spaceGUID"`
	MemoryMB  int64    `json:"memoryMB"`
	DiskMB    int64    `json:"diskMB"`
	// +kubebuilder:validation:Format:=uint8
	CPUWeight uint8 `json:"cpuWeight"`
}

type ExecutionStatus string

const (
	TaskStarting  ExecutionStatus = "starting"
	TaskRunning   ExecutionStatus = "running"
	TaskSucceeded ExecutionStatus = "succeeded"
	TaskFailed    ExecutionStatus = "failed"
)

type TaskStatus struct {
	StartTime *metav1.Time `json:"start_time"`
	EndTime   *metav1.Time `json:"end_time"`
	// +kubebuilder:validation:Enum=starting;running;succeeded;failed
	// +kubebuilder:default=starting
	ExecutionStatus ExecutionStatus `json:"execution_status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Task is the Schema for the tasks API
type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskSpec   `json:"spec,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TaskList contains a list of Task
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Task{}, &TaskList{})
}
