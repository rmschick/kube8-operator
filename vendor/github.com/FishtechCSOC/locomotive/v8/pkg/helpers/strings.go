package helpers

func NilToString(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

func DefaultString(defaultString, checkedString string) string {
	switch {
	case checkedString != "":
		return checkedString
	default:
		return defaultString
	}
}
