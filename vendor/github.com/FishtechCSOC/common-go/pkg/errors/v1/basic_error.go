package errors

var _ error = BasicError("")

type BasicError string

func (err BasicError) Error() string {
	return string(err)
}
