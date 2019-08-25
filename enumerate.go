package enumerate

import (
	"bytes"
	"errors"
	"go/format"
	"io"
	"text/template"
)

// Write writes an enum file to an io.Writer.
func Write(e *Enum, w io.Writer) error {
	if err := validate(e); err != nil {
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

func validate(enum *Enum) error {
	if enum.Package == "" {
		return errors.New("package required")
	}

	if enum.Type == "" {
		return errors.New("type required")
	}

	return nil
}
