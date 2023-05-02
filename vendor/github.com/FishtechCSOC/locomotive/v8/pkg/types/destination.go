package types

import (
	"strings"
)

type Destination string

const (
	Splunk                 Destination = "splunk"
	Chronicle              Destination = "chronicle"
	ChronicleUDM           Destination = "chronicle-udm"
	ChronicleEntities      Destination = "chronicle-entities"
	Azure                  Destination = "azure"
	Elastic                Destination = "elastic"
	Exabeam                Destination = "exabeam"
	LogzIO                 Destination = "logz-io"
	SumoLogic              Destination = "sumologic"
	QRadar                 Destination = "qradar"
	Alerting               Destination = "alerting"
	Gatekeeper             Destination = "gatekeeper"
	LongTermStorage        Destination = "lts"
	InsiderThreatDetection Destination = "itd"
	GCS                    Destination = "gcs"
	// Deprecated: no longer applicable.
	Haystax Destination = "haystax"
	// Deprecated: no longer applicable.
	CDPAlerting Destination = "cdp-alerting"
	// Deprecated: no longer applicable.
	CDPDetection Destination = "cdp-detection"
	// Deprecated: no longer applicable.
	DataSourceMetrics Destination = "data-source-metrics"

	destinationHeaderPrefix = "x-cyderes-destination-"

	defaultSeparator = ","
)

func UnmarshalDestinationsFromMetadata(metadata string) []Destination {
	if metadata == "" {
		return []Destination{}
	}

	list := strings.Split(metadata, defaultSeparator)

	destinations := make([]Destination, len(list))

	for i, value := range list {
		destinations[i] = Destination(strings.ToLower(value))
	}

	return destinations
}

func UnmarshalDestinationMapFromMetadata(metahash map[string]string) []Destination {
	var destinations []Destination

	for key := range metahash {
		if strings.Contains(key, destinationHeaderPrefix) {
			destinations = append(destinations, Destination(
				strings.ToLower(
					strings.ReplaceAll(key, destinationHeaderPrefix, ""),
				),
			),
			)

			delete(metahash, key)
		}
	}

	delete(metahash, destinationKey)

	return destinations
}

func MarshalDestinationsToString(destinations ...Destination) string {
	list := make([]string, len(destinations))

	for i, destination := range destinations {
		list[i] = string(destination)
	}

	return strings.Join(list, defaultSeparator)
}

func MarshalDestinationsToMap(destinations ...Destination) map[string]string {
	destinationMap := make(map[string]string)

	for _, destination := range destinations {
		destinationMap[destinationHeaderPrefix+string(destination)] = "true"
	}

	return destinationMap
}

// Deprecated: in favor of `ContainsDestination`.
func Contains(destinations []Destination, destination Destination) bool {
	return ContainsDestination(destinations, destination)
}

func ContainsDestination(sources []Destination, target Destination) bool {
	for _, value := range sources {
		if value == target {
			return true
		}
	}

	return false
}

func ContainsAnyDestination(sources []Destination, targets []Destination) bool {
	for _, value := range sources {
		if ContainsDestination(targets, value) {
			return true
		}
	}

	return false
}

func DefaultDestinations(defaultDataType, checkedDataType []Destination) []Destination {
	switch {
	case len(checkedDataType) > 0:
		return checkedDataType
	default:
		return defaultDataType
	}
}
