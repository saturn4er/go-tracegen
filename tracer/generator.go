package tracer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type Generator struct {
	Dir          string
	OutputFile   string
	Prefix       string
	OnlyFiles    []string
	ExcludeFiles []string
	OnlyTypes    []string
	ExcludeTypes []string
}

func (g *Generator) Generate() error {
	types, err := g.parse()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := g.render(types); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (g *Generator) parse() ([]Type, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, g.Dir, nil, 0)
	if err != nil {
		panic(err)
	}
	types := make(map[string]Type)
	for pkgNm, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			if !g.needToParseFile(fileName) {
				continue
			}

			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range decl.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						interfaceSpec, ok := typeSpec.Type.(*ast.InterfaceType)
						if !ok {
							continue
						}
						interfaceName := typeSpec.Name.String()

						use, err := g.needToUseType(interfaceName)
						if err != nil {
							return nil, err
						}
						if !use {
							continue
						}
						interfaceType := Type{
							PkgName:     pkgNm,
							Name:        interfaceName,
							IsInterface: true,
						}

						for _, method := range interfaceSpec.Methods.List {
							funcType, ok := method.Type.(*ast.FuncType)
							if !ok {
								continue
							}
							interfaceType.Methods = append(interfaceType.Methods, funcTypeToMethod(method.Names[0].Name, funcType))
						}
						types[typeSpec.Name.String()] = interfaceType
					}
				case *ast.FuncDecl:
					if decl.Recv == nil {
						continue
					}
					var recvType string
					if startExpr, ok := decl.Recv.List[0].Type.(*ast.StarExpr); ok {
						recvType, err = exprToString(startExpr.X)
						if err != nil {
							return nil, err
						}
					} else {
						recvType, err = exprToString(decl.Recv.List[0].Type)
						if err != nil {
							return nil, err
						}
					}

					use, err := g.needToUseType(recvType)
					if err != nil {
						return nil, err
					}
					if !use {
						continue
					}

					if rune, _ := utf8.DecodeRuneInString(decl.Name.String()); unicode.IsLower(rune) {
						continue
					}
					typ, ok := types[recvType]
					if !ok {
						typ = Type{
							PkgName: pkgNm,
							Name:    recvType,
						}
					}

					typ.Methods = append(typ.Methods, funcTypeToMethod(decl.Name.String(), decl.Type))
					types[recvType] = typ
				}
			}

		}
	}
	res := make([]Type, 0, len(types))
	for _, typ := range types {
		res = append(res, typ)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})

	return res, nil
}

func (g *Generator) render(types []Type) error {
	if len(types) == 0 {
		return nil
	}

	tp := template.New("tpl")
	tp, err := tp.Parse(tpl)
	if err != nil {
		return errors.WithStack(err)
	}
	file, err := os.OpenFile(g.OutputFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	pkgName := types[0].PkgName
	err = tp.Execute(file, map[string]interface{}{
		"pkgName": pkgName,
		"tracers": types,
		"prefix":  g.Prefix,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (g *Generator) needToParseFile(filename string) bool {
	if len(g.OnlyFiles) > 0 {
		var found bool
		for _, file := range g.OnlyFiles {
			if file == filename {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	for _, excludedFilename := range g.ExcludeFiles {
		if excludedFilename == filename {
			return false
		}
	}

	return true
}

func (g *Generator) needToUseType(typeName string) (bool, error) {
	excludeGlobs, err := stringsToGobs(g.ExcludeTypes)
	if err != nil {
		return false, err
	}

	onlyGlobs, err := stringsToGobs(g.OnlyTypes)
	if err != nil {
		return false, err
	}

	if len(onlyGlobs) > 0 {
		var found bool
		for _, glob := range onlyGlobs {
			if glob.Match(typeName) {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	for _, glob := range excludeGlobs {
		if glob.Match(typeName) {
			return false, nil
		}
	}

	return true, nil
}

func NewGenerator(dir, outputFile, prefix string, onlyFiles, excludeFiles, onlyTypes, excludeTypes []string) *Generator {
	return &Generator{
		Dir:          dir,
		OutputFile:   outputFile,
		Prefix:       prefix,
		OnlyFiles:    onlyFiles,
		ExcludeFiles: excludeFiles,
		OnlyTypes:    onlyTypes,
		ExcludeTypes: excludeTypes,
	}
}
