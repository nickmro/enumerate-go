package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"os"
	"strings"

	"github.com/nickmro/gonumerate"
)

const usageText = `Usage:
 gonumerate [OPTION]...

Options:
 -type        The enum type name (Required)
 -description The type comment description (must have format "<type> is...")
 -values      The enum values
 -prefix      The prefix to apply to each enum value
 -json        The JSON encoding type {string, int}
 -sql         The SQL encoding type {string, int}
 -o           The output filename
 -help        Print usage
`

var (
	typeOpt   = flag.String("type", "", "The enum type name (Required)")
	descOpt   = flag.String("description", "", "The type comment description (must be format \"<type> is...\")")
	valuesOpt = flag.String("values", "", "The comma-separated list of enum values")
	prefixOpt = flag.String("prefix", "", "The prefix to apply to each enum value")
	jsonOpt   = flag.String("json", "", "The JSON encoding type {string, int}")
	sqlOpt    = flag.String("sql", "", "The SQL encoding type {string, int}")
	outputOpt = flag.String("o", "", "The output filename")
	printOpt  = flag.Bool("help", false, "Print usage")
)

func main() {
	var e gonumerate.Enum
	var err error

	if flag.Parse(); !flag.Parsed() {
		fmt.Println(usageText)
		os.Exit(1)
	}

	if printOpt != nil && *printOpt {
		fmt.Println(usageText)
		os.Exit(1)
	}

	e.Type, err = parseType()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	e.JSONEncoding, err = parseJSON()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	e.SQLEncoding, err = parseSQL()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	e.Package, err = parsePackageName()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	e.Description = parseDescription()
	e.Values = parseValues()
	e.Prefix = parsePrefix()
	e.OutFile = praseOutFile()

	f, err := os.Create(e.FileName())
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	err = e.Write(w)
	if err != nil {
		printError(err)
		os.Remove(e.FileName())
		os.Exit(1)
	}
}

func parseType() (string, error) {
	if typeOpt == nil || *typeOpt == "" {
		return "", errors.New("type required")
	}
	return *typeOpt, nil
}

func parseDescription() string {
	if descOpt != nil && *descOpt != "" {
		return *descOpt
	}
	return ""
}

func parseValues() []string {
	if valuesOpt != nil && *valuesOpt != "" {
		return strings.Split(*valuesOpt, ",")
	}
	return []string{}
}

func parsePrefix() string {
	if prefixOpt != nil {
		return *prefixOpt
	}
	return ""
}

func parseJSON() (gonumerate.Encoding, error) {
	if jsonOpt != nil && *jsonOpt != "" {
		if j := gonumerate.EncodingFromString(*jsonOpt); j != 0 {
			return j, nil
		} else {
			return 0, errors.New("invalid json value")
		}
	}
	return 0, nil
}

func parseSQL() (gonumerate.Encoding, error) {
	if sqlOpt != nil && *sqlOpt != "" {
		if s := gonumerate.EncodingFromString(*sqlOpt); s != 0 {
			return s, nil
		} else {
			return 0, errors.New("invalid sql value")
		}
	}
	return 0, nil
}

func praseOutFile() string {
	if outputOpt != nil && *outputOpt != "" {
		return *outputOpt
	}
	return ""
}

func parsePackageName() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pkg, err := build.ImportDir(wd, 0)
	if err != nil {
		return "", err
	}

	return pkg.Name, nil
}

func printError(e error) {
	fmt.Println("Error:", e)
	fmt.Println()
	fmt.Println(usageText)
}
