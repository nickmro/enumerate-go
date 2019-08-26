package gonumerate_test

import (
	"bytes"
	"testing"

	"github.com/nickmro/gonumerate"
)

func TestWrite(t *testing.T) {
	e := enum()
	b := bytes.NewBuffer([]byte{})
	err := e.Write(b)
	if err != nil {
		t.Error("expected", err, "to be nil")
	}
}

func TestFileName(t *testing.T) {
	e := enum()
	n := e.FileName()
	if n != "user_type.go" {
		t.Error("expected", n, "to equal", "user_type.go")
	}
}

func TestImports(t *testing.T) {
	e := enum()

	// When JSON and SQL are specified
	i := e.Imports()
	if len(i) != 3 {
		t.Error("expected", i, "to have length", 3)
	}
	if i[0] != "encoding/json" {
		t.Error("expected", i[0], "to equal", "encoding/json")
	}
	if i[1] != "database/sql/driver" {
		t.Error("expected", i[1], "to equal", "database/sql/driver")
	}
	if i[2] != "errors" {
		t.Error("expected", i[2], "to equal", "errors")
	}

	// When SQL is not specified
	e = enum()
	e.SQLEncoding = 0
	i = e.Imports()
	if len(i) != 1 {
		t.Error("expected", i, "to have length", 1)
	}
	if i[0] != "encoding/json" {
		t.Error("expected", i[0], "to equal", "encoding/json")
	}

	// When JSON is not specified
	e = enum()
	e.JSONEncoding = 0
	i = e.Imports()
	if len(i) != 2 {
		t.Error("expected", i, "to have length", 2)
	}
	if i[0] != "database/sql/driver" {
		t.Error("expected", i[0], "to equal", "database/sql/driver")
	}
	if i[1] != "errors" {
		t.Error("expected", i[1], "to equal", "errors")
	}

	// When no encoding is specified
	e = enum()
	e.JSONEncoding = 0
	e.SQLEncoding = 0
	i = e.Imports()
	if len(i) != 0 {
		t.Error("expected", i, "to have length", 0)
	}
}

func TestStringMapName(t *testing.T) {
	e := enum()
	n := e.StringMapName()
	if n != "userTypeStrings" {
		t.Error("expected", n, "to equal", "userTypeStrings")
	}
}

func TestConstructorName(t *testing.T) {
	e := enum()
	n := e.ConstructorName()
	if n != "UserTypeFromString" {
		t.Error("expected", n, "to equal", "UserTypeFromString")
	}
}

func TestValueNames(t *testing.T) {
	e := enum()
	n := e.ValueNames()
	if len(n) != 2 {
		t.Fatal("expected", n, "to have length", 2)
		t.FailNow()
	}
	if n[0] != "UserTypeAdmin" {
		t.Error("expected", n[0], "to equal", "UserTypeAdmin")
	}
	if n[1] != "UserTypeCustomerSupport" {
		t.Error("expected", n[1], "to equal", "UserTypeCustomerSupport")
	}
}

func TestMappedStrings(t *testing.T) {
	e := enum()
	s := e.MappedStrings()

	expected := []string{
		"UserTypeAdmin: \"admin\"",
		"UserTypeCustomerSupport: \"customer_support\"",
	}

	if len(s) != 2 {
		t.Error("expected", s, "to have length", 2)
	}

	for i, v := range expected {
		if s[i] != v {
			t.Error("expected", s[i], "to equal", v)
		}
	}
}

func TestValidate(t *testing.T) {
	e := enum()

	// When enum is complete
	err := e.Validate()
	if err != nil {
		t.Error("expected", err, "to be nil")
	}

	// When package is not specified
	e = enum()
	e.Package = ""
	err = e.Validate()
	if err == nil {
		t.Error("expected", err, "to occur")
	}
	if err != gonumerate.ErrPackageRequred {
		t.Error("expected", err, "to equal", gonumerate.ErrPackageRequred)
	}

	// When type is not specified
	e = enum()
	e.Type = ""
	err = e.Validate()
	if err == nil {
		t.Error("expected", err, "to occur")
	}
	if err != gonumerate.ErrTypeRequired {
		t.Error("expected", err, "to equal", gonumerate.ErrTypeRequired)
	}
}

func enum() *gonumerate.Enum {
	return &gonumerate.Enum{
		Package:      "package_name",
		Type:         "UserType",
		Prefix:       "UserType",
		Values:       []string{"Admin", "CustomerSupport"},
		JSONEncoding: gonumerate.EncodingString,
		SQLEncoding:  gonumerate.EncodingString,
	}
}
