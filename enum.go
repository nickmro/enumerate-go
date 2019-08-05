package enumerate

// Enum represents an enumeration file template.
type Enum struct {
	Package      string   // The package name
	Type         string   // The enum type name
	Values       []string // The enum values
	Prefix       string   // The prefix to apply to each enum
	JSONEncoding Encoding // The JSON encoding type
}

const enumTemplate = `
package {{.Package}}

import (
	{{- if .JSONEncoding}}
	"encoding/json"
	{{- end}}
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

{{- if .JSONEncoding}}
// MarshalJSON marshals the {{.Type}} to JSON.
func (t {{.Type}}) MarshalJSON() ([]byte, error) {
	{{- if eq .JSONEncoding "string" }}
	return json.Marshal(t.String())
	{{- else}}
	return t
	{{- end}}
}

// UnmarshalJSON unmarshals the {{.Type}} from JSON.
func (t *{{.Type}}) UnmarshalJSON(b []byte) error {
	{{- if eq .JSONEncoding "string" }}
	var v string
	{{- else}}
	var v int
	{{- end}}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	{{- if eq .JSONEncoding "string" }}
	*t = {{ toCamelCase .Type }}FromString(v)
	{{- else}}
	*t = v
	{{- end}}
}
{{- end}}

func {{ toCamelCase .Type }}FromString(s string) {{.Type}} {
	for k, v := range {{ toCamelCase .Type }}Strings {
		if v == s {
			return k
		}
	}
	return 0
}

`
