package gonumerate

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"
)

// Enum represents an enumeration file template.
type Enum struct {
	Package      string   // The package name
	Type         string   // The enum type name
	Description  string   // The enum description
	Values       []string // The enum values
	Prefix       string   // The prefix to apply to each enum
	JSONEncoding Encoding // The JSON encoding type
	SQLEncoding  Encoding // The SQL encoding type
	OutFile      string   // The output filename
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

{{- if .Description}}
// {{.Description}}
{{- else }}
// {{.Type}} is an enumeration of values
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

// Enum errors
var (
	ErrPackageRequred     = errors.New("package required")
	ErrTypeRequired       = errors.New("type required")
	ErrDescriptionInvalid = errors.New("description invalid")
)

// Write writes an enum file to an io.Writer.
func (e *Enum) Write(w io.Writer) error {
	if err := e.Validate(); err != nil {
		return err
	}

	// Create a new template
	t, err := template.New("enum").
		Parse(enumTemplate)
	if err != nil {
		return err
	}

	// Write the template to a buffer
	var buf bytes.Buffer
	err = t.Execute(&buf, e)
	if err != nil {
		return err
	}

	// Format the code
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	// Write the code to a file
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// FileName returns the enum's file name.
func (e Enum) FileName() string {
	if e.OutFile != "" {
		return e.OutFile
	}

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
		names = append(names, valueName(e.Prefix, v))
	}
	return names
}

// MappedStrings returns the strings mapped to a value.
func (e Enum) MappedStrings() []string {
	s := []string{}
	for _, v := range e.Values {
		b := strings.Builder{}
		b.WriteString(valueName(e.Prefix, v))
		b.WriteString(": \"")
		b.WriteString(valueString((v)))
		b.WriteString("\"")
		s = append(s, b.String())
	}
	return s
}

// Validate returns an error if the enum is invalid.
func (e *Enum) Validate() error {
	if e.Package == "" {
		return ErrPackageRequred
	}

	if e.Type == "" {
		return ErrTypeRequired
	}

	if e.Description != "" && !strings.HasPrefix(e.Description, fmt.Sprintf("%s is", e.Type)) {
		return ErrDescriptionInvalid
	}

	return nil
}

func valueName(prefix, value string) string {
	b := strings.Builder{}
	b.WriteString(toPascalCase(prefix))
	b.WriteString(toPascalCase(value))
	return b.String()
}

func valueString(v string) string {
	return toSnakeCase(v)
}
