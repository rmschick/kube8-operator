package configuration

import (
	"strings"
)

const (
	separator = "."
)

type Defaults struct {
	prefix   string
	fields   map[string]any
	children []*Defaults
}

func CreateDefaults(prefix string) *Defaults {
	return &Defaults{
		prefix:   prefix,
		fields:   make(map[string]any),
		children: make([]*Defaults, 0),
	}
}

func (config *Defaults) Copy() *Defaults {
	newCtx := &Defaults{
		prefix:   config.prefix,
		fields:   make(map[string]any),
		children: make([]*Defaults, len(config.children)),
	}

	for k, v := range config.fields {
		newCtx.fields[k] = v
	}

	for i, child := range config.children {
		newCtx.children[i] = child.Copy()
	}

	return newCtx
}

func (config *Defaults) WithPrefix(prefix string) *Defaults {
	newCtx := config.Copy()

	newCtx.prefix = prefix

	return newCtx
}

func (config *Defaults) WithFields(fields map[string]any) *Defaults {
	newCtx := config.Copy()

	for k, v := range fields {
		newCtx.fields[k] = v
	}

	return newCtx
}

func (config *Defaults) WithChildren(children ...*Defaults) *Defaults {
	newCtx := config.Copy()

	newCtx.children = append(newCtx.children, children...)

	return newCtx
}

func (config *Defaults) GetFields() map[string]any {
	return config.getFlattenedFields([]string{})
}

func (config *Defaults) getFlattenedFields(paths []string) map[string]any {
	fields := make(map[string]any)

	if config.prefix != "" {
		paths = append(paths, config.prefix)
	}

	for k, v := range config.fields {
		fields[strings.Join(append(paths, k), separator)] = v
	}

	return config.gatherChildFields(paths, fields)
}

func (config *Defaults) gatherChildFields(paths []string, fields map[string]any) map[string]any {
	for _, child := range config.children {
		results := child.getFlattenedFields(paths)

		for k, v := range results {
			fields[k] = v
		}
	}

	return fields
}
