package tracer

import (
	"go/ast"
	"strings"

	"github.com/gobwas/glob"
)

func stringsToGobs(strs []string) ([]glob.Glob, error) {
	res := make([]glob.Glob, 0, len(strs))
	for _, str := range strs {
		glob, err := glob.Compile(str)
		if err != nil {
			return nil, err
		}
		res = append(res, glob)
	}

	return res, nil
}

func exprToString(typ ast.Expr) (string, error) {
	switch typ := typ.(type) {
	case *ast.MapType:
		keyTyp, err := exprToString(typ.Key)
		if err != nil {
			return "", err
		}
		valueTyp, err := exprToString(typ.Value)
		if err != nil {
			return "", err
		}
		return "map[" + keyTyp + "]" + valueTyp, nil
	case *ast.ArrayType:
		itemTyp, err := exprToString(typ.Elt)
		if err != nil {
			return "", err
		}
		arrLen, err := exprToString(typ.Len)
		if err != nil {
			return "", err
		}
		return "[" + arrLen + "]" + itemTyp, nil
	case *ast.InterfaceType:
		return "interface{}", nil
	case *ast.StarExpr:
		elemTyp, err := exprToString(typ.X)
		if err != nil {
			return "", err
		}
		return "*" + elemTyp, nil
	case *ast.FuncType:
		funStr := "func"
		funStr += "(" + fieldListToStr(typ.Params) + ")"
		if typ.Results != nil {
			funStr += " (" + fieldListToStr(typ.Results) + ")"
		}

		return funStr, nil
	case *ast.SelectorExpr:
		xTyp, err := exprToString(typ.X)
		if err != nil {
			return "", err
		}
		valTyp, err := exprToString(typ.Sel)
		if err != nil {
			return "", err
		}
		return xTyp + "." + valTyp, nil
	case *ast.Ellipsis:
		elType, err := exprToString(typ.Elt)
		if err != nil {
			return "", err
		}
		return "..." + elType, nil
	case *ast.Ident:
		return typ.Name, nil
	case nil:
		return "", nil
	default:
		panic(typ)
	}

}

func fieldListToStr(fieldList *ast.FieldList) string {
	if fieldList == nil {
		return ""
	}
	var params []string
	for _, param := range fieldList.List {
		typName, err := exprToString(param.Type)
		if err != nil {
			panic(err)
		}
		var names []string
		for _, name := range param.Names {
			nameVal, err := exprToString(name)
			if err != nil {
				panic(err)
			}
			names = append(names, nameVal)
		}
		if len(names) == 0 {
			params = append(params, typName)
		} else {
			params = append(params, strings.Join(names, ",")+" "+typName)
		}
	}

	return strings.Join(params, ", ")
}

func funcTypeToMethod(name string, funcType *ast.FuncType) TypeMethod {
	method := TypeMethod{
		Name: name,
	}
	for _, param := range funcType.Params.List {
		typeStr, err := exprToString(param.Type)
		if err != nil {
			panic(err)
		}
		if len(param.Names) != 0 {
			for _, name := range param.Names {
				_, isEllipsis := param.Type.(*ast.Ellipsis)
				method.Params = append(method.Params, MethodParam{
					Name:       name.Name,
					Type:       typeStr,
					IsEllipsis: isEllipsis,
				})
			}
		} else {
			method.Params = append(method.Params, MethodParam{
				Name: "",
				Type: typeStr,
			})
		}
	}
	if funcType.Results != nil {
		i := 0
		for _, result := range funcType.Results.List {
			typeStr, err := exprToString(result.Type)
			if err != nil {
				panic(err)
			}
			if len(result.Names) != 0 {
				for _, name := range result.Names {
					method.Results = append(method.Results, MethodParam{
						Name:  name.Name,
						Type:  typeStr,
						Index: i,
					})
					i++
				}
			} else {
				method.Results = append(method.Results, MethodParam{
					Name:  "",
					Type:  typeStr,
					Index: i,
				})
				i++
			}
		}
	}

	return method
}
