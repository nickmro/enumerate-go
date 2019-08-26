# Gonumerate
[![Build Status](https://travis-ci.com/nickmro/gonumerate.svg?branch=master)](https://travis-ci.com/nickmro/gonumerate)

Gonumerate is an enumeration generator for Go.

## Why?

Writing enumerations in Go can be achieved using a type alias, but adding methods to represent those enumerations as strings can be reptitive. This code generator eases the process.

## Installation

```
go get -u github.com/nickmro/gonumerate
cd $GOPATH/src/github.com/nickmro/gonumerate/cmd/gonumerate
go build -o $GOPATH/bin/gonumerate
```

## Usage

```
Usage:
 gonumerate [OPTION]...

Options:
 -type    The enum type name (Required)
 -values  The enum values
 -prefix  The prefix to apply to each enum value
 -json    The JSON encoding type {string, int}
 -sql     The SQL encoding type {string, int}
 -o       The output filename
 -help    Print usage
```

## Example

```bash
gonumerate --type UserType \
	--prefix UserType \
	--values Admin,Support \
	--json string \
	--sql string
```

This will produce the following file:
```go
package gonumerate

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

// String returns a string representation of the UserType.
func (t UserType) String() string {
	if v, ok := userTypeStrings[t]; ok {
		return v
	}
	return ""
}

// UserTypeFromString returns the UserType from the given string.
func UserTypeFromString(s string) UserType {
	for k, v := range userTypeStrings {
		if v == s {
			return k
		}
	}
	return 0
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
	*t = UserTypeFromString(v)
	return nil
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
		*t = UserTypeFromString(string(b))
		return nil
	} else if s, ok := bv.(string); ok {
		*t = UserTypeFromString(s)
		return nil
	} else {
		*t = 0
		return errors.New("failed to scan UserType")
	}
}
```

## go generate

This command can be run with `go generate`. Simply add the command as a comment to any `.go` file in your directory. For example:

```go
//go:generate gonumerate --type UserType --values Admin,Support
```

Then run:
```bash
go generate
```

For more information about `go generate`: https://blog.golang.org/generate

## Other options

This is not the only enumeration generator written in Go. For other options, see:

- https://github.com/abice/go-enum
- https://github.com/alvaroloes/enumer
- https://github.com/steinfletcher/gonum
