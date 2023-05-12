package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupName string = "example.com"
	Kind      string = "Service"
	Version   string = "v1"
	Plural    string = "services"
	Singular  string = "service"
	ShortName string = "svc"
	Name      string = Plural + "." + GroupName
)

type ServiceSpec struct {
	Type     string            `json:"type"`
	Selector map[string]string `json:"selector"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Service describes a Service resource
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec `json:"spec,omitempty"`
	Status string      `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceList is a list of Service resources
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Service `json:"items"`
}
