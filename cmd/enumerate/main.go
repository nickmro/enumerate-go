package main

import (
	"bufio"
	"flag"
	"go/build"
	"os"
	"strings"

	"github.com/nickmro/enumerate-go"
)

func main() {
	var e enumerate.Enum

	enumType := flag.String("type", "", "The enum type name")
	enumValues := flag.String("values", "", "The comma-separated list of enum values")
	enumPrefix := flag.String("prefix", "", "The prefix to apply to each enum value")

	flag.Parse()

	if t := enumType; t != nil {
		e.Type = *t
	}

	if v := enumValues; v != nil {
		e.Values = strings.Split(*v, ",")
	}

	if p := enumPrefix; p != nil {
		e.Prefix = *p
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p, err := getPackageName(wd)
	if err != nil {
		panic(err)
	}

	e.Package = p

	f, err := os.Create(fileName(e.Type))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	err = enumerate.Write(&e, w)
	if err != nil {
		panic(err)
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
