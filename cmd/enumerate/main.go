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

const usageText = `Usage:
 enumerate [OPTION]...

Options:
 -type    The enum type name (Required)
 -values  The enum values
 -prefix  The prefix to apply to each enum value
 -json    The JSON encoding type {string, int}
 -sql     The SQL encoding type {string, int}
 -help    Print usage
`

func main() {
	typeOpt := flag.String("type", "", "The enum type name (Required)")
	valuesOpt := flag.String("values", "", "The comma-separated list of enum values")
	prefixOpt := flag.String("prefix", "", "The prefix to apply to each enum value")
	jsonOpt := flag.String("json", "", "The JSON encoding type {string, int}")
	sqlOpt := flag.String("sql", "", "The SQL encoding type {string, int}")
	printOpt := flag.Bool("help", false, "Print usage")

	var e enumerate.Enum

	if flag.Parse(); !flag.Parsed() {
		fmt.Println(usageText)
		os.Exit(1)
	}

	if printOpt != nil && *printOpt {
		fmt.Println(usageText)
		os.Exit(1)
	}

	if typeOpt == nil || *typeOpt == "" {
		fmt.Println(usageText)
		os.Exit(1)
	}

	e.Type = *typeOpt

	if valuesOpt != nil && *valuesOpt != "" {
		e.Values = strings.Split(*valuesOpt, ",")
	}

	if prefixOpt != nil {
		e.Prefix = *prefixOpt
	}

	if jsonOpt != nil && *jsonOpt != "" {
		if j := enumerate.EncodingFromString(*jsonOpt); j != 0 {
			e.JSONEncoding = j
		} else {
			fmt.Println(usageText)
			os.Exit(1)
		}
	}

	if sqlOpt != nil && *sqlOpt != "" {
		if s := enumerate.EncodingFromString(*sqlOpt); s != 0 {
			e.SQLEncoding = s
		} else {
			fmt.Println(usageText)
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
