package views

import (
	"errors"
	"fmt"
	"html/template"
	"neft.web/errorController"
)

type Template struct {
	Layout string
	File   string
	View   *View
}

var Templates map[string]Template

func InitTemplateController() {
	Templates = make(map[string]Template)

	Templates["LiroMain"] = Template{
		Layout: "dashboard",
		File:   "liro/index",
		View:   NewReload("dashboard", "views/liro/index.gohtml"),
	}
	Templates["LiroUserList"] = Template{
		Layout: "dashboard",
		File:   "liro/listUsers",
		View:   NewReload("dashboard", "views/liro/listUsers.gohtml"),
	}
}

func ReloadHtml() {

	for key, template := range Templates {
		template.View.Template = reload(template.File)
		errorController.DebugLogger.Println("new view reloaded:", key)
	}
	errorController.DebugLogger.Println("all views reloaded")
}

func reload(files ...string) *template.Template {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, LayoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrf is not implemented")
		},
	}).ParseFiles(files...)
	if err != nil {
		fmt.Println("err", err)
	}
	return t

}
