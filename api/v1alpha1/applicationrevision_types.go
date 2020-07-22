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

// ApplicationRevisionSpec defines the desired state of ApplicationRevision
type ApplicationRevisionSpec struct {
	App     Application       `json:"app,omitempty"`
	Cfg     ApplicationConfig `json:"cfg,omitempty"`
	AppHash string            `json:"appHash,omitempty"`
	CfgHash string            `json:"cfgHash,omitempty"`
}

// ApplicationRevisionStatus defines the observed state of ApplicationRevision
type ApplicationRevisionStatus struct {
}

// +kubebuilder:object:root=true

// ApplicationRevision is the Schema for the applicationrevisions API
type ApplicationRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationRevisionSpec   `json:"spec,omitempty"`
	Status ApplicationRevisionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationRevisionList contains a list of ApplicationRevision
type ApplicationRevisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationRevision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationRevision{}, &ApplicationRevisionList{})
}
