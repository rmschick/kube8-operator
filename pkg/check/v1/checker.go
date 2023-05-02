package check

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/check/v1"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ check.Checker = (*Checker)(nil)

type Checker struct{}

func NewChecker() Checker {
	return Checker{}
}

func (checker Checker) Info() string {
	return `TODO`
}

func (checker Checker) Fields() []check.Field {
	return []check.Field{}
}

func (checker Checker) Test(_ context.Context, data map[string]any) (any, error) {
	if err := validateFields(data); err != nil {
		return nil, err
	}

	var result any

	return result, nil
}

func validateFields(data map[string]any) error {
	return validation.Validate(data,
		validation.Map(),
	)
}
