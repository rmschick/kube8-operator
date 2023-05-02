package ecs

import (
	"bytes"
	"encoding/json"
)

// JSONMarshalNoHTMLEscape disables unicode-escapes for &, <, and >
// which may be present inside URLs.
func JSONMarshalNoHTMLEscape(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)

	return buffer.Bytes(), err
}
