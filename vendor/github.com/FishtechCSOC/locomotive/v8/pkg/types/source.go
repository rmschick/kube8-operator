package types

import (
	"strings"

	"github.com/FishtechCSOC/locomotive/v8/pkg/helpers"
)

type Source struct {
	Type           string            `json:"type" mapstructure:"type"`
	Path           string            `json:"path" mapstructure:"path"`
	Agent          string            `json:"agent" mapstructure:"agent"`
	Infrastructure map[string]string `json:"infrastructure" mapstructure:"infrastructure"`
}

func (source *Source) DeepCopy() Source {
	deepCopy := Source{
		Infrastructure: make(map[string]string),
		Type:           source.Type,
		Path:           source.Path,
		Agent:          source.Agent,
	}

	for k, v := range source.Infrastructure {
		deepCopy.Infrastructure[k] = v
	}

	return deepCopy
}

func (source *Source) MergeLeft(other Source) Source {
	return Source{
		Agent:          helpers.DefaultString(other.Agent, source.Agent),
		Path:           helpers.DefaultString(other.Path, source.Path),
		Type:           helpers.DefaultString(other.Type, source.Type),
		Infrastructure: helpers.MergeStringMapLeft(other.Infrastructure, source.Infrastructure),
	}
}

func (source *Source) MergeRight(other Source) Source {
	return Source{
		Agent:          helpers.DefaultString(source.Agent, other.Agent),
		Path:           helpers.DefaultString(source.Path, other.Path),
		Type:           helpers.DefaultString(source.Type, other.Type),
		Infrastructure: helpers.MergeStringMapRight(other.Infrastructure, source.Infrastructure),
	}
}

func UnmarshalInfrastructureMapFromMetadata(metahash map[string]string) map[string]string {
	infrastructureLabels := make(map[string]string)

	for key, value := range metahash {
		if !strings.Contains(key, sourceInfrastructureKey) {
			continue
		}

		infrastructureLabels[strings.ToLower(strings.ReplaceAll(key, sourceInfrastructureKey+"-", ""))] = value

		delete(metahash, key)
	}

	return infrastructureLabels
}

func (source *Source) Empty() bool {
	return source.Type == "" && source.Path == "" && source.Agent == "" && len(source.Infrastructure) <= 0
}
