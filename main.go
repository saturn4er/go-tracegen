package main

import (
	"fmt"
	"os"

	"github.com/saturn4er/go-tracegen/tracer"
	"github.com/spf13/pflag"
)

type MethodParam struct {
	Name       string
	Type       string
	IsEllipsis bool
	Index      int
}

func (m MethodParam) IsError() bool {
	return m.Type == "error"
}

type TypeMethod struct {
	Name    string
	Params  []MethodParam
	Results []MethodParam
}

func (t TypeMethod) ErrorsResults() []MethodParam {
	var res []MethodParam
	for _, result := range t.Results {
		if result.IsError() {
			res = append(res, result)
		}
	}

	return res
}

func (t TypeMethod) HaveErrorResults() bool {
	for _, result := range t.Results {
		if result.IsError() {
			return true
		}
	}

	return false
}

func (t TypeMethod) NeedTracing() bool {
	for _, param := range t.Params {
		if param.Type == "context.Context" {
			return true
		}
	}

	return false
}

func (t TypeMethod) NeedReturn() bool {
	return len(t.Results) > 0
}

type Type struct {
	Name        string
	IsInterface bool
	Methods     []TypeMethod
}

func (t Type) NextType() string {
	if t.IsInterface {
		return t.Name
	}

	return "*" + t.Name
}

func main() {
	var (
		prefix       = pflag.String("prefix", "", "Span names prefix")
		dir          = pflag.String("dir", ".", "Where to search for types")
		out          = pflag.String("out", "tracer.go", "Output file")
		onlyTypes    = pflag.StringArray("only-types", []string{}, "Which types to use")
		excludeTypes = pflag.StringArray("ignore-types", []string{}, "Which types to ignore")
		onlyFiles    = pflag.StringArray("only-files", []string{}, "Which files to parse")
		excludeFiles = pflag.StringArray("ignore-files", []string{}, "Which files to ignore")
	)
	pflag.Parse()

	generator := tracer.NewGenerator(*dir, *out, *prefix, *onlyFiles, *excludeFiles, *onlyTypes, *excludeTypes)

	if err := generator.Generate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
