package types

import (
	"fmt"
)

// Enums are sadly a seriously depressing lack of syntactic sugar in golang, so this looks more awful than it should be.
type DataType string

const (
	ChronicleDefaultMapping = "DEFAULT"
	AzureDefaultMapping     = "CommonSecurityLog"
)

// Deprecated: We should use the CMS to validate data types if we do it at all.
func (logType DataType) MapToChronicle() string {
	if _, ok := chronicleMapping[logType]; !ok {
		return ChronicleDefaultMapping
	}

	return chronicleMapping[logType]
}

// Deprecated: We should use the CMS to validate data types if we do it at all.
func (logType DataType) MustMapToChronicle() (string, error) {
	if _, ok := chronicleMapping[logType]; !ok {
		return "", fmt.Errorf("unrecognized log type for Chronicle: %s", logType)
	}

	return chronicleMapping[logType], nil
}

// Deprecated: We should use the CMS to validate data types if we do it at all.
func (logType DataType) MapToAzure() string {
	if _, ok := azureMapping[logType]; !ok {
		return AzureDefaultMapping
	}

	return azureMapping[logType]
}

// Deprecated: We should use the CMS to validate data types if we do it at all.
func (logType DataType) MustMapToAzure() (string, error) {
	if _, ok := azureMapping[logType]; !ok {
		return "", fmt.Errorf("unrecognized log type for Azure: %s", logType)
	}

	return azureMapping[logType], nil
}

// Deprecated: We should use the CMS to validate data types if we do it at all.
func DefaultDataType(defaultDataType, checkedDataType DataType) DataType {
	switch {
	case checkedDataType != "":
		return checkedDataType
	default:
		return defaultDataType
	}
}
