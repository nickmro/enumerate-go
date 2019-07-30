package enumerate

import (
	"bytes"
	"errors"
	"go/format"
	"io"
	"regexp"
	"strings"
	"text/template"
)

var (
	matchFirstCap      = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAnyCap        = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	matchFirstChar     = regexp.MustCompile(`(^[A-Za-z])`)
	matchFirstWordChar = regexp.MustCompile(`_([A-Za-z])`)
)

// Write writes an enum file to an io.Writer.
func Write(e *Enum, w io.Writer) error {
	if err := validate(e); err != nil {
		return err
	}

	f := template.FuncMap{
		"toCamelCase": ToCamelCase,
		"toSnakeCase": ToSnakeCase,
	}

	// Create a new template
	t, err := template.New("enum").
		Funcs(f).
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

func ToSnakeCase(v string) string {
	v = matchFirstCap.ReplaceAllString(v, `${1}_${2}`)
	v = matchAnyCap.ReplaceAllString(v, `${1}_${2}`)
	return strings.ToLower(v)
}

func ToCamelCase(v string) string {
	v = matchFirstWordChar.ReplaceAllStringFunc(v, func(s string) string {
		return strings.ToUpper(strings.Replace(s, `_`, ``, -1))
	})

	// lowercase first character
	r := []rune(v)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]

	return string(r)
}

func validate(enum *Enum) error {
	if enum.Package == "" {
		return errors.New("package required")
	}

	if enum.Type == "" {
		return errors.New("type required")
	}

	return nil
}
