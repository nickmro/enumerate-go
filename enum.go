package enumerate

// Enum represents an enumeration file template.
type Enum struct {
	Package string   // The package name
	Type    string   // The enum type name
	Values  []string // The enum values
	Prefix  string   // The prefix to apply to each enum
}

const enumTemplate = `
package {{.Package}}

import (
	"encoding/json"
)

type {{.Type}} int

// The {{.Type}} values
const (
	_ {{.Type}} = iota
	{{- $prefix := .Prefix}}
	{{- range $i, $v := .Values}}
	{{$prefix}}{{$v}}
	{{- end}}
)

var {{ toCamelCase .Type }}Strings = map[{{.Type}}]string{
	{{- $prefix := .Prefix}}
	{{- range $i, $v := .Values}}
	{{$prefix}}{{$v}}: "{{ toSnakeCase $v }}",
	{{- end}}
}

// String returns a string representation of the {{.Type}}
func (t {{.Type}}) String() string  {
	if v, ok := {{ toCamelCase .Type }}Strings[t]; ok {
		return v
	}
	return ""
}

// MarshalJSON marshals the {{.Type}} to JSON.
func (t {{.Type}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON unmarshals the {{.Type}} from JSON.
func (t *{{.Type}}) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	for k, v := range {{ toCamelCase .Type }}Strings {
		if v == s {
			*t = k
		}
	}
	return nil
}
`
