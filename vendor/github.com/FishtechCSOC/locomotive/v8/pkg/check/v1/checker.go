package check

import "context"

type Checker interface {
	Info() string
	Fields() []Field
	Test(context.Context, map[string]any) (any, error)
}
