package main

type TypeTracerTracer struct {
	next *TypeTracer
}

func (t *TypeTracerTracer) NextType() (_ string) {
	return t.next.NextType()
}

func NewTypeTracerTracer(next *TypeTracer) *TypeTracerTracer {
	return &TypeTracerTracer{
		next: next,
	}
}

type TypeMethodTracerTracer struct {
	next *TypeMethodTracer
}

func (t *TypeMethodTracerTracer) ErrorsResults() (_ []MethodParam) {
	return t.next.ErrorsResults()
}
func (t *TypeMethodTracerTracer) HaveErrorResults() (_ bool) {
	return t.next.HaveErrorResults()
}
func (t *TypeMethodTracerTracer) NeedTracing() (_ bool) {
	return t.next.NeedTracing()
}
func (t *TypeMethodTracerTracer) NeedReturn() (_ bool) {
	return t.next.NeedReturn()
}

func NewTypeMethodTracerTracer(next *TypeMethodTracer) *TypeMethodTracerTracer {
	return &TypeMethodTracerTracer{
		next: next,
	}
}

type TypeMethodTracer struct {
	next *TypeMethod
}

func (t *TypeMethodTracer) ErrorsResults() (_ []MethodParam) {
	return t.next.ErrorsResults()
}
func (t *TypeMethodTracer) HaveErrorResults() (_ bool) {
	return t.next.HaveErrorResults()
}
func (t *TypeMethodTracer) NeedTracing() (_ bool) {
	return t.next.NeedTracing()
}
func (t *TypeMethodTracer) NeedReturn() (_ bool) {
	return t.next.NeedReturn()
}

func NewTypeMethodTracer(next *TypeMethod) *TypeMethodTracer {
	return &TypeMethodTracer{
		next: next,
	}
}

type TypeTracer struct {
	next *Type
}

func (t *TypeTracer) NextType() (_ string) {
	return t.next.NextType()
}

func NewTypeTracer(next *Type) *TypeTracer {
	return &TypeTracer{
		next: next,
	}
}

type MethodParamTracerTracer struct {
	next *MethodParamTracer
}

func (t *MethodParamTracerTracer) IsError() (_ bool) {
	return t.next.IsError()
}

func NewMethodParamTracerTracer(next *MethodParamTracer) *MethodParamTracerTracer {
	return &MethodParamTracerTracer{
		next: next,
	}
}

type MethodParamTracer struct {
	next *MethodParam
}

func (t *MethodParamTracer) IsError() (_ bool) {
	return t.next.IsError()
}

func NewMethodParamTracer(next *MethodParam) *MethodParamTracer {
	return &MethodParamTracer{
		next: next,
	}
}
