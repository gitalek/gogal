package views

import "html/template"

func NewView(layout string, files ...string) (*View, error) {
	files = append(
		files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml",
	)
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
