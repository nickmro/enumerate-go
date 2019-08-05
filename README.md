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
enumerate -type=${TYPE_NAME} -values=${TYPE_VALUES} -prefix=${TYPE_PREFIX} -json=${JSON_ENCODING}
```

```
-json string
    The JSON encoding type [string, int]
-prefix string
    The prefix to apply to each enum value
-type string
    The enum type name
-values string
    The comma-separated list of enum values
```

## Example

```bash
enumerate -type=UserType -values=Admin,Support -prefix=UserType -json=string
```

This will produce the type:
```go
const UserType int
```

with the values:
```go
const (
    _ UserType = iota
    UserTypeAdmin
    UserTypeSupport
)
```

## Other options

This is not the only enumeration generator written in Go. For other options, see:

- https://github.com/abice/go-enum
- https://github.com/alvaroloes/enumer
- https://github.com/steinfletcher/gonum
