package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"os"
	"strings"

	"github.com/nickmro/enumerate-go"
)

func main() {
	enumType := flag.String("type", "", "The enum type name (Required)")
	enumValues := flag.String("values", "", "The comma-separated list of enum values")
	enumPrefix := flag.String("prefix", "", "The prefix to apply to each enum value")
	enumJSON := flag.String("json", "", "The JSON encoding type {string, int}")
	enumSQL := flag.String("sql", "", "The SQL encoding type {string, int}")

	var e enumerate.Enum

	if flag.Parse(); !flag.Parsed() {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if enumType == nil || *enumType == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	e.Type = *enumType

	if enumValues != nil && *enumValues != "" {
		e.Values = strings.Split(*enumValues, ",")
	}

	if enumPrefix != nil {
		e.Prefix = *enumPrefix
	}

	if enumJSON != nil && *enumJSON != "" {
		if j := enumerate.EncodingFromString(*enumJSON); j != 0 {
			e.JSONEncoding = j
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if enumSQL != nil && *enumSQL != "" {
		if s := enumerate.EncodingFromString(*enumSQL); s != 0 {
			e.SQLEncoding = s
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p, err := getPackageName(wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.Package = p

	f, err := os.Create(fileName(e.Type))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	err = enumerate.Write(&e, w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getPackageName(dir string) (string, error) {
	pkg, err := build.ImportDir(dir, 0)
	if err != nil {
		return "", err
	}

	return pkg.Name, nil
}

func fileName(name string) string {
	b := strings.Builder{}
	b.WriteString(enumerate.ToSnakeCase(name))
	b.WriteString(".go")
	return b.String()
}
