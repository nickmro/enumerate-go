package enumerate

import (
	"strings"
)

// Enum represents an enumeration file template.
type Enum struct {
	Package      string   // The package name
	Type         string   // The enum type name
	Values       []string // The enum values
	Prefix       string   // The prefix to apply to each enum
	JSONEncoding Encoding // The JSON encoding type
	SQLEncoding  Encoding // The SQL encoding type
}

const enumTemplate = `package {{.Package}}

{{- $imports := .Imports}}
{{- if gt (len $imports) 0}}
import (
	{{- range $i, $import := $imports}}
	"{{ $import }}"
	{{- end}}
)
{{- end}}

type {{.Type}} int

// The {{.Type}} values
const (
	_ {{.Type}} = iota
	{{- range $i, $n := .ValueNames}}
	{{ $n }}
	{{- end}}
)

var {{ .StringMapName }} = map[{{.Type}}]string{
	{{- range $i, $s := .MappedStrings}}
	{{ $s }},
	{{- end}}
}

// String returns a string representation of the {{.Type}}.
func (t {{.Type}}) String() string  {
	if v, ok := {{ .StringMapName }}[t]; ok {
		return v
	}
	return ""
}

// {{ .ConstructorName }} returns the {{.Type}} from the given string.
func {{ .ConstructorName }}(s string) {{.Type}} {
	for k, v := range {{ .StringMapName }} {
		if v == s {
			return k
		}
	}
	return 0
}

{{- if .JSONEncoding}}
// MarshalJSON marshals the {{.Type}} to JSON.
func (t {{.Type}}) MarshalJSON() ([]byte, error) {
	{{- if eq .JSONEncoding 1 }}
	return json.Marshal(t.String())
	{{- else}}
	return t
	{{- end}}
}

// UnmarshalJSON unmarshals the {{.Type}} from JSON.
func (t *{{.Type}}) UnmarshalJSON(b []byte) error {
	{{- if eq .JSONEncoding 1 }}
	var v string
	{{- else}}
	var v int
	{{- end}}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	{{- if eq .JSONEncoding 1 }}
	*t = {{ .ConstructorName }}(v)
	{{- else}}
	*t = v
	{{- end}}
	return nil
}
{{- end}}

{{- if .SQLEncoding}}
// Value returns the {{.Type}} value for SQL encoding.
func (t *{{.Type}}) Value() (driver.Value, error) {
	{{- if eq .SQLEncoding 1 }}
	return t.String(), nil
	{{- else}}
	return t, nil
	{{- end}}
}

// Scan scans the {{.Type}} from its SQL encoded value.
func (t *{{.Type}}) Scan(v interface{}) error {
	{{- if eq .SQLEncoding 1 }}
	bv, err := driver.String.ConvertValue(v)
	if err != nil {
		*t = 0
		return errors.New("failed to scan {{.Type}}")
	}

	if b, ok := bv.([]byte); ok {
		*t = {{ .ConstructorName }}(string(b))
		return nil
	} else if s, ok := bv.(string); ok {
		*t = {{ .ConstructorName }}(s)
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

`

// FileName returns the enum's file name.
func (e Enum) FileName() string {
	b := strings.Builder{}
	b.WriteString(toSnakeCase(e.Type))
	b.WriteString(".go")
	return b.String()
}

// Imports returns an enum file's imports.
func (e Enum) Imports() []string {
	imports := []string{}
	if e.JSONEncoding != 0 {
		imports = append(imports, "encoding/json")
	}
	if e.SQLEncoding != 0 {
		imports = append(imports,
			"database/sql/driver",
			"errors")
	}
	return imports
}

// StringMapName returns the name of the string map.
func (e Enum) StringMapName() string {
	b := strings.Builder{}
	b.WriteString(toCamelCase(e.Type))
	b.WriteString("Strings")
	return b.String()
}

// ConstructorName returns the name of the constructor function.
func (e Enum) ConstructorName() string {
	b := strings.Builder{}
	b.WriteString(toPascalCase(e.Type))
	b.WriteString("FromString")
	return b.String()
}

// ValueNames returns the value names.
func (e Enum) ValueNames() []string {
	names := []string{}
	for _, v := range e.Values {
		b := strings.Builder{}
		b.WriteString(toPascalCase(e.Prefix))
		b.WriteString(toPascalCase(v))
		names = append(names, b.String())
	}
	return names
}

// MappedStrings returns the strings mapped to a value.
func (e Enum) MappedStrings() []string {
	s := []string{}
	for _, v := range e.Values {
		b := strings.Builder{}
		b.WriteString(toPascalCase(e.Prefix))
		b.WriteString(toPascalCase(v))
		b.WriteString(": \"")
		b.WriteString(toSnakeCase(v))
		b.WriteString("\"")
		s = append(s, b.String())
	}
	return s
}
