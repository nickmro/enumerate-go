package enumerate

// Enum represents an enumeration file template.
type Enum struct {
	Package      string   // The package name
	Type         string   // The enum type name
	Values       []string // The enum values
	Prefix       string   // The prefix to apply to each enum
	JSONEncoding Encoding // The JSON encoding type
	SQLEncoding  Encoding // The SQL encoding type
}

const enumTemplate = `
package {{.Package}}

import (
	{{- if .JSONEncoding}}
	"encoding/json"
	{{- end}}
	{{- if .SQLEncoding}}
	"database/sql/driver"
	"errors"
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

{{- if .SQLEncoding}}
// Value returns the {{.Type}} value for SQL encoding.
func (t *{{.Type}}) Value() (driver.Value, error) {
	{{- if eq .SQLEncoding "string" }}
	return t.String(), nil
	{{- else}}
	return t, nil
	{{- end}}
}

// Scan scans the {{.Type}} from its SQL encoded value.
func (t *{{.Type}}) Scan(v interface{}) error {
	{{- if eq .SQLEncoding "string" }}
	bv, err := driver.String.ConvertValue(v)
	if err != nil {
		*t = 0
		return errors.New("failed to scan {{.Type}}")
	}

	if b, ok := bv.([]byte); ok {
		*t = {{toCamelCase .Type}}FromString(string(b))
		return nil
	} else if s, ok := bv.(string); ok {
		*t = {{toCamelCase .Type}}FromString(s)
		return nil
	} else {
		*t = 0
		return errors.New("failed to scan {{.Type}}")
	}
	{{- else}}
	if b, ok := v.(int); ok {
		*t = {{.Type}}(b)
		return nil
	}
	return errors.New("failed to scan {{.Type}}")
	{{- end}}
}
{{- end}}

func {{toCamelCase .Type}}FromString(s string) {{.Type}} {
	for k, v := range {{ toCamelCase .Type }}Strings {
		if v == s {
			return k
		}
	}
	return 0
}

`
