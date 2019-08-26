package main

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap      = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAnyCap        = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	matchFirstChar     = regexp.MustCompile(`(^[A-Za-z])`)
	matchFirstWordChar = regexp.MustCompile(`_([A-Za-z])`)
)

func toSnakeCase(v string) string {
	v = matchFirstCap.ReplaceAllString(v, `${1}_${2}`)
	v = matchAnyCap.ReplaceAllString(v, `${1}_${2}`)
	return strings.ToLower(v)
}

func toCamelCase(v string) string {
	v = matchFirstWordChar.ReplaceAllStringFunc(v, func(s string) string {
		return strings.ToUpper(strings.Replace(s, `_`, ``, -1))
	})

	// lowercase first character
	r := []rune(v)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]

	return string(r)
}

func toPascalCase(v string) string {
	v = matchFirstWordChar.ReplaceAllStringFunc(v, func(s string) string {
		return strings.ToUpper(strings.Replace(s, `_`, ``, -1))
	})

	// uppercase first character
	r := []rune(v)
	r[0] = []rune(strings.ToUpper(string(r[0])))[0]
	return string(r)
}
