package templates

import (
	"html/template"
	"io"
	"path/filepath"
)

type Templates struct {
	sources  []string
	template *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		sources:  []string{},
		template: nil,
	}
}

func (t *Templates) AddSource(source string) {
	t.sources = append(t.sources, source)
}

func (t *Templates) Parse() error {
	if t.template == nil {
		t.template = template.New("main")
	}
	for _, source := range t.sources {
		path := filepath.Join(source, "*.html")
		_, err := t.template.ParseGlob(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Templates) Execute(wr io.Writer, name string, data interface{}) error {
	return t.template.ExecuteTemplate(wr, name, data)
}
