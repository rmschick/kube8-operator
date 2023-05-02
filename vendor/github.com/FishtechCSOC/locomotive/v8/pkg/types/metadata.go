package types

import (
	"strings"

	"github.com/FishtechCSOC/locomotive/v8/pkg/helpers"
)

const (
	CyderesKeyPrefix        = "x-cyderes"
	clientIDKey             = CyderesKeyPrefix + "-client-id"
	clientNameKey           = CyderesKeyPrefix + "-client-name"
	dataTypeKey             = CyderesKeyPrefix + "-data-type"
	destinationKey          = CyderesKeyPrefix + "-destination"
	environmentKey          = CyderesKeyPrefix + "-environment"
	timestampFieldKey       = CyderesKeyPrefix + "-timestamp-field"
	persistentObjectKey     = CyderesKeyPrefix + "-persistent-object"
	batchIDKey              = CyderesKeyPrefix + "-batch-id"
	sourceAgentKey          = CyderesKeyPrefix + "-source-agent"
	sourcePathKey           = CyderesKeyPrefix + "-source-path"
	sourceTypeKey           = CyderesKeyPrefix + "-source-type"
	sourceInfrastructureKey = CyderesKeyPrefix + "-source-infrastructure"
	namespaceKey            = CyderesKeyPrefix + "-namespace"
	tenantIDKey             = CyderesKeyPrefix + "-tenant-id"
	tenantReferenceKey      = CyderesKeyPrefix + "-tenant-reference"
)

// Metadata is an object that gets attached to logs and contains info
// about destination, customer, etc.
// nolint: tagliatelle
type Metadata struct {
	Labels       map[string]string `json:"labels" mapstructure:"labels"`
	BatchID      string            `json:"batchID" mapstructure:"batchID"`
	Customer     Customer          `json:"customer" mapstructure:"customer"`
	Environment  string            `json:"environment" mapstructure:"environment"`
	Source       Source            `json:"source" mapstructure:"source"`
	DataType     DataType          `json:"dataType" mapstructure:"dataType"`
	Destinations []Destination     `json:"destinations" mapstructure:"destinations"`
	// Deprecated
	TimestampField   string `json:"timestampField" mapstructure:"timestampField"`
	PersistentObject string `json:"persistentObject" mapstructure:"persistentObject"`
	Instance         string `json:"instance" mapstructure:"instance"`
	Namespace        string `json:"namespace" mapstructure:"namespace"`
	Tenant           Tenant `json:"tenant" mapstructure:"tenant"`
}

func CreateMetadata() Metadata {
	return Metadata{
		Labels:       make(map[string]string),
		Destinations: make([]Destination, 0),
		Source: Source{
			Agent: CreateUserAgent(),
		},
	}
}

func (metadata *Metadata) Get(key string) string {
	if metadata.Labels == nil {
		return ""
	}

	value, ok := metadata.Labels[key]
	if !ok {
		return ""
	}

	return value
}

func (metadata *Metadata) Set(key string, value string) {
	if metadata.Labels == nil {
		metadata.Labels = make(map[string]string)
	}

	metadata.Labels[key] = value
}

func (metadata *Metadata) Keys() []string {
	keys := make([]string, 0, len(metadata.Labels))

	for k := range metadata.Labels {
		keys = append(keys, k)
	}

	return keys
}

func (metadata *Metadata) MarshalToMetadata() map[string]string {
	metaHash := make(map[string]string)

	for k, v := range metadata.Labels {
		if k == "" {
			continue
		}

		metaHash[k] = v
	}

	for k, v := range metadata.Source.Infrastructure {
		if k == "" {
			continue
		}

		metaHash[sourceInfrastructureKey+"-"+k] = v
	}

	for k, v := range MarshalDestinationsToMap(metadata.Destinations...) {
		metaHash[k] = v
	}

	metaHash[clientIDKey] = metadata.Customer.ID
	metaHash[clientNameKey] = metadata.Customer.Name
	metaHash[dataTypeKey] = string(metadata.DataType)
	metaHash[environmentKey] = metadata.Environment
	metaHash[destinationKey] = MarshalDestinationsToString(metadata.Destinations...)
	metaHash[batchIDKey] = metadata.BatchID
	metaHash[sourceAgentKey] = metadata.Source.Agent
	metaHash[sourcePathKey] = metadata.Source.Path
	metaHash[sourceTypeKey] = metadata.Source.Type
	metaHash[timestampFieldKey] = metadata.TimestampField
	metaHash[persistentObjectKey] = metadata.PersistentObject
	metaHash[namespaceKey] = metadata.Namespace
	metaHash[tenantIDKey] = metadata.Tenant.ID
	metaHash[tenantReferenceKey] = metadata.Tenant.Reference

	return metaHash
}

func (metadata *Metadata) UnmarshalFromMetadata(hash map[string]string) {
	metaHash := copyHash(hash)

	switch {
	case checkForPrefix(metaHash):
		metadata.Destinations = UnmarshalDestinationMapFromMetadata(metaHash)
	default:
		metadata.Destinations = UnmarshalDestinationsFromMetadata(popValue(metaHash, destinationKey))
	}

	metadata.Customer.ID = popValue(metaHash, clientIDKey)
	metadata.DataType = DataType(popValue(metaHash, dataTypeKey))
	metadata.Environment = popValue(metaHash, environmentKey)
	metadata.Customer.Name = popValue(metaHash, clientNameKey)
	metadata.BatchID = popValue(metaHash, batchIDKey)
	metadata.Source.Agent = popValue(metaHash, sourceAgentKey)
	metadata.Source.Path = popValue(metaHash, sourcePathKey)
	metadata.Source.Type = popValue(metaHash, sourceTypeKey)
	metadata.Source.Infrastructure = UnmarshalInfrastructureMapFromMetadata(metaHash)
	metadata.TimestampField = popValue(metaHash, timestampFieldKey)
	metadata.PersistentObject = popValue(metaHash, persistentObjectKey)
	metadata.Namespace = popValue(metaHash, namespaceKey)
	metadata.Tenant.ID = popValue(metaHash, tenantIDKey)
	metadata.Tenant.Reference = popValue(metaHash, tenantReferenceKey)

	delete(metaHash, "")

	metadata.Labels = metaHash
}

// nolint: dupl
func (metadata *Metadata) MergeLeft(other Metadata) Metadata {
	return Metadata{
		DataType:         DefaultDataType(other.DataType, metadata.DataType),
		BatchID:          helpers.DefaultString(other.BatchID, metadata.BatchID),
		Environment:      helpers.DefaultString(other.Environment, metadata.Environment),
		PersistentObject: helpers.DefaultString(other.PersistentObject, metadata.PersistentObject),
		TimestampField:   helpers.DefaultString(other.TimestampField, metadata.TimestampField),
		Customer:         metadata.Customer.MergeLeft(other.Customer),
		Source:           metadata.Source.MergeLeft(other.Source),
		Destinations:     DefaultDestinations(other.Destinations, metadata.Destinations),
		Labels:           helpers.MergeStringMapLeft(metadata.Labels, other.Labels),
		Instance:         helpers.DefaultString(other.Instance, metadata.Instance),
		Namespace:        helpers.DefaultString(other.Namespace, metadata.Namespace),
		Tenant:           metadata.Tenant.MergeLeft(other.Tenant),
	}
}

// nolint: dupl
func (metadata *Metadata) MergeRight(other Metadata) Metadata {
	return Metadata{
		DataType:         DefaultDataType(metadata.DataType, other.DataType),
		BatchID:          helpers.DefaultString(metadata.BatchID, other.BatchID),
		Environment:      helpers.DefaultString(metadata.Environment, other.Environment),
		PersistentObject: helpers.DefaultString(metadata.PersistentObject, other.PersistentObject),
		TimestampField:   helpers.DefaultString(metadata.TimestampField, other.TimestampField),
		Customer:         metadata.Customer.MergeRight(other.Customer),
		Source:           metadata.Source.MergeRight(other.Source),
		Destinations:     DefaultDestinations(metadata.Destinations, other.Destinations),
		Labels:           helpers.MergeStringMapRight(metadata.Labels, other.Labels),
		Instance:         helpers.DefaultString(metadata.Instance, other.Instance),
		Namespace:        helpers.DefaultString(metadata.Namespace, other.Namespace),
		Tenant:           metadata.Tenant.MergeRight(other.Tenant),
	}
}

// nolint: gocyclo, cyclop
// Empty checks if metadata struct is populated.
func (metadata *Metadata) Empty() bool {
	if metadata.DataType != "" {
		return false
	}

	if len(metadata.Destinations) > 0 {
		return false
	}

	if !metadata.Customer.Empty() {
		return false
	}

	if metadata.Environment != "" {
		return false
	}

	if !metadata.Source.Empty() {
		return false
	}

	if metadata.BatchID != "" {
		return false
	}

	if metadata.TimestampField != "" {
		return false
	}

	if metadata.PersistentObject != "" {
		return false
	}

	if len(metadata.Labels) > 0 {
		return false
	}

	if metadata.Instance != "" {
		return false
	}

	if metadata.Namespace != "" {
		return false
	}

	if metadata.Tenant.Empty() {
		return false
	}

	return true
}

func copyHash(src map[string]string) map[string]string {
	dst := make(map[string]string)

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func popValue(hash map[string]string, key string) string {
	defer delete(hash, key)

	return hash[key]
}

func checkForPrefix(hash map[string]string) bool {
	for k := range hash {
		if strings.Contains(k, destinationHeaderPrefix) {
			return true
		}
	}

	return false
}
