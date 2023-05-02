package check

const (
	TEXT        FieldType = "TEXT"
	SELECT      FieldType = "SELECT"
	MULTISELECT FieldType = "MULTISELECT"
)

type FieldType string

type Field struct {
	FieldType    FieldType `json:"fieldType"`
	Name         string    `json:"name"`
	Display      string    `json:"display"`
	Hint         string    `json:"hint"`
	DefaultValue string    `json:"defaultValue"`
	IsRequired   bool      `json:"isRequired"`
	IsSecure     bool      `json:"isSecure"`
	IsMultiLine  bool      `json:"isMultiline"`
	Options      []string  `json:"options"`
}
