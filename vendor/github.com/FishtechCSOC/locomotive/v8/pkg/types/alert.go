package types

import (
	"github.com/FishtechCSOC/ecs-golang/pkg/ecs"
	"github.com/FishtechCSOC/ecs-golang/pkg/elasticsiem"
)

// Deprecated: this has been migrated to github.com/FishtechCSOC/morter.
type Alert struct {
	Cyderes Cyderes  // cyderes fields
	ECS     ecs.Base // ecs format
	SIEM    elasticsiem.ElasticSIEM
}

// Deprecated: this has been migrated to github.com/FishtechCSOC/morter.
type Cyderes struct {
	Name        string      `json:"name,omitempty" mapstructure:"name"`               // alertname
	Description string      `json:"description,omitempty" mapstructure:"description"` // Whatever doesn't fit into a Field will be dropped into here
	IssueType   IssueType   `json:"issuetype,omitempty" mapstructure:"issuetype"`     // The type of alert/issue, under Project
	Project     ProjectType `json:"project,omitempty" mapstructure:"project"`         // Project determines who works on it
	ClientID    string      `json:"@clientRef,omitempty" mapstructure:"@clientRef"`   // Client ID
	ClientName  string      `json:"@clientName,omitempty" mapstructure:"@clientName"` // Client Name
}

// Deprecated: this has been migrated to github.com/FishtechCSOC/morter.
type IssueType string

const (
	THREATHUNT  IssueType = "THREATHUNT"
	ALERT       IssueType = "ALERT"
	IOC         IssueType = "IOC"
	HEALTHALERT IssueType = "HEALTHALERT"
)

// Deprecated: this has been migrated to github.com/FishtechCSOC/morter.
type ProjectType string

const (
	ACTION ProjectType = "ACTION"
	CSOC   ProjectType = "CSOC"
	CASE   ProjectType = "CASE"
)
