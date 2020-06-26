package tracer

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
	PkgName     string
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
