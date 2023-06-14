package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupName string = "example.com"
	Kind      string = "Collector"
	Version   string = "v1"
	Plural    string = "collectors"
	Singular  string = "collector"
	ShortName string = "collector"
	Name      string = Plural + "." + GroupName
)

type CollectorInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	Configuration string `json:"configuration"`
}

type TenantInfo struct {
	ID        string `json:"id"`
	Reference string `json:"reference"`
	Instance  string `json:"instance"`
}

type CollectorSpec struct {
	Collector CollectorInfo `json:"collector"`
	Tenant    TenantInfo    `json:"tenant"`
	Cluster   string        `json:"cluster"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

type Collector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CollectorSpec   `json:"spec,omitempty"`
	Status CollectorStatus `json:"status,omitempty"`
}

// CollectorStatus defines the observed state of Collector
type CollectorStatus struct {
	Status metav1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CollectorList is a list of Collector resources
type CollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Collector `json:"items"`
}
