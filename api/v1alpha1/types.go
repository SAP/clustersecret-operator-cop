/*
Copyright 2023 SAP SE.

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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sap/component-operator-runtime/pkg/component"
	componentoperatorruntimetypes "github.com/sap/component-operator-runtime/pkg/types"
)

// ClusterSecretOperatorSpec defines the desired state of ClusterSecretOperator
type ClusterSecretOperatorSpec struct {
	component.Spec `json:",inline"`
	// +optional
	Controller ControllerSpec `json:"controller"`
	// +optional
	Webhook WebhookSpec `json:"webhook"`
}

// ControllerSpec defines the desired state of the controller deployment
type ControllerSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	ReplicaCount int `json:"replicaCount,omitempty"`
	// +optional
	Image                          component.ImageSpec `json:"image"`
	component.KubernetesProperties `json:",inline"`
	LogLevel                       int `json:"logLevel,omitempty"`
}

// WebhookSpec defines the desired state of the webhook deployment
type WebhookSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	ReplicaCount int `json:"replicaCount,omitempty"`
	// +optional
	Image                          component.ImageSpec `json:"image"`
	component.KubernetesProperties `json:",inline"`
	LogLevel                       int `json:"logLevel,omitempty"`
}

// ClusterSecretOperatorStatus defines the observed state of ClusterSecretOperator
type ClusterSecretOperatorStatus struct {
	component.Status `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterSecretOperator is the Schema for the clustersecretoperators API
type ClusterSecretOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterSecretOperatorSpec `json:"spec,omitempty"`
	// +kubebuilder:default={"observedGeneration":-1}
	Status ClusterSecretOperatorStatus `json:"status,omitempty"`
}

var _ component.Component = &ClusterSecretOperator{}

// +kubebuilder:object:root=true

// ClusterSecretOperatorList contains a list of ClusterSecretOperator
type ClusterSecretOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterSecretOperator `json:"items"`
}

func (s *ClusterSecretOperatorSpec) ToUnstructured() map[string]any {
	result, err := runtime.DefaultUnstructuredConverter.ToUnstructured(s)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ClusterSecretOperator) GetDeploymentNamespace() string {
	if c.Spec.Namespace != "" {
		return c.Spec.Namespace
	}
	return c.Namespace
}

func (c *ClusterSecretOperator) GetDeploymentName() string {
	if c.Spec.Name != "" {
		return c.Spec.Name
	}
	return c.Name
}

func (c *ClusterSecretOperator) GetSpec() componentoperatorruntimetypes.Unstructurable {
	return &c.Spec
}

func (c *ClusterSecretOperator) GetStatus() *component.Status {
	return &c.Status.Status
}

func PostReadHook(ctx context.Context, client client.Client, c *ClusterSecretOperator) error {
	if c.Spec.Namespace == "" {
		c.Spec.Namespace = c.Namespace
	}
	if c.Spec.Name == "" {
		c.Spec.Name = c.Name
	}
	return nil
}

func init() {
	SchemeBuilder.Register(&ClusterSecretOperator{}, &ClusterSecretOperatorList{})
}
