package views

import "html/template"

func NewView(files ...string) (*View, error) {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return &View{Template: t}, nil
}

type View struct {
	Template *template.Template
}
