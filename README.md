# Enumerate

Enumerate is an enumeration generator for Go.

## Why?

Writing enumerations in Go is often a very repetitive task, which is an ideal case for a code generator.

## Installation

```bash
go get -u github.com/nickmro/enumerate-go
cd $GOPATH/src/github.com/nickmro/enumerate-go/cmd/enumerate
go build -o $GOPATH/bin/enumerate
```

## Instructions

In the directory of the package for which you would like to create the enumeration:

```bash
enumerate -type=${TYPE_NAME} -values=${TYPE_VALUES} -prefix=${TYPE_PREFIX} -json=${JSON_ENCODING} -sql=${SQL_ENCODING}
```

```
-json string
    The JSON encoding type [string, int]
-prefix string
    The prefix to apply to each enum value
-sql string
    The SQL encoding type [string, int]
-type string
    The enum type name
-values string
    The comma-separated list of enum values
```

## Example

```bash
enumerate -type=UserType -values=Admin,Support -prefix=UserType -json=string -sql=string
```

This will produce the following file:
```go
package enumerate

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type UserType int

// The UserType values
const (
	_ UserType = iota
	UserTypeAdmin
	UserTypeSupport
)

var userTypeStrings = map[UserType]string{
	UserTypeAdmin:   "admin",
	UserTypeSupport: "support",
}

// String returns a string representation of the UserType
func (t UserType) String() string {
	if v, ok := userTypeStrings[t]; ok {
		return v
	}
	return ""
}

// MarshalJSON marshals the UserType to JSON.
func (t UserType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON unmarshals the UserType from JSON.
func (t *UserType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*t = userTypeFromString(v)
}

// Value returns the UserType value for SQL encoding.
func (t *UserType) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan scans the UserType from its SQL encoded value.
func (t *UserType) Scan(v interface{}) error {
	bv, err := driver.String.ConvertValue(v)
	if err != nil {
		*t = 0
		return errors.New("failed to scan UserType")
	}

	if b, ok := bv.([]byte); ok {
		*t = userTypeFromString(string(b))
		return nil
	} else if s, ok := bv.(string); ok {
		*t = userTypeFromString(s)
		return nil
	} else {
		*t = 0
		return errors.New("failed to scan UserType")
	}
}

func userTypeFromString(s string) UserType {
	for k, v := range userTypeStrings {
		if v == s {
			return k
		}
	}
	return 0
}
```

## Other options

This is not the only enumeration generator written in Go. For other options, see:

- https://github.com/abice/go-enum
- https://github.com/alvaroloes/enumer
- https://github.com/steinfletcher/gonum
