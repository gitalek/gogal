package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) (*View, error) {
	layoutFiles, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		return nil, err
	}
	files = append(files, layoutFiles...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return &View{Template: t, Layout: layout}, nil
}

type View struct {
	Template *template.Template
	Layout   string
}
