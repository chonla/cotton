package template

import (
	"cotton/internal/variable"

	"github.com/valyala/fasttemplate"
)

type Template struct {
	tmplRaw    string
	tmplEngine *fasttemplate.Template
}

func New(tmpl string) *Template {
	return &Template{
		tmplRaw:    tmpl,
		tmplEngine: fasttemplate.New(tmpl, "{{", "}}"),
	}
}

func (t *Template) Apply(variables *variable.Variables) string {
	if variables == nil {
		return t.tmplRaw
	}
	values := variables.ToStringMap()
	return t.tmplEngine.ExecuteStringStd(values)
}
