package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"neft.web/errorController"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func newView(layout string, files ...string) *View {

	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, LayoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrf is not implemented")
		},
	}).ParseFiles(files...)
	if err != nil {
		errorController.ErrorLogger.Println(err)
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Parsing templates"
		errorController.WD.SendErrorWHWeb()
		return nil
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-type", "text/html")
	var vd Data

	switch d := data.(type) {

	case *Data:
		vd = *d
	default:
		vd.Yield = d
	}
	vd.Active = r.URL.Path
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)

	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		errorController.ErrorLogger.Println(err)
		http.Redirect(w, r, "/505", http.StatusFound)
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Error executing template"
		errorController.WD.SendErrorWHWeb()
		return
	}
	_, err := io.Copy(w, &buf)
	if err != nil {
		errorController.ErrorLogger.Println(err)
		http.Redirect(w, r, "/505", http.StatusFound)
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Error executing template"
		errorController.WD.SendErrorWHWeb()
		return
	}
}

func LayoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Error generating template files"
		errorController.WD.SendErrorWHWeb()
		return nil
	}
	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
