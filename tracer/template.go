package tracer

const tpl = `package {{$.pkgName -}}
{{range $tracer := $.tracers }}
type {{$tracer.Name -}}Tracer struct{
	next {{$tracer.NextType}}
}
	
{{range $i, $method := $tracer.Methods}}
func (t *{{$tracer.Name -}}Tracer) {{$method.Name}}({{ range $param := $method.Params -}}{{$param.Name}} {{$param.Type}},{{ end -}})({{ range $result := $method.Results -}}{{if $result.IsError}}rerr{{$result.Index}}{{else}}_{{end}} {{$result.Type}},{{ end -}}){ 
	{{ if $method.NeedTracing -}}
	ctx, span := tracing.ChildSpanFromContext(ctx, "{{$.prefix}}{{$tracer.Name}}.Get")
	{{ if $method.HaveErrorResults }}
	defer func() { tracing.FinishSpanWithErr(span, multierr.Combine({{range $err := $method.ErrorsResults}}rerr{{$err.Index}},{{end}})) }()
	{{ else }}
	defer span.Finish()
	{{ end }}
	
	{{ end -}}
	{{if $method.NeedReturn }}return{{end}} t.next.{{$method.Name}}({{range $param := $method.Params}}{{$param.Name}}{{if $param.IsEllipsis}}...{{end}},{{end}})
}
{{- end}}
	
func New{{$tracer.Name -}}Tracer(next {{$tracer.NextType -}}) *{{$tracer.Name -}}Tracer{
	return &{{$tracer.Name -}}Tracer{
		next: next,
	}
}

{{end}}
`
