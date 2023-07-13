package controllers

import (
	"neft.web/views"
	"net/http"
)

func NewStatic() *Static {
	return &Static{
		Home:     views.CreateView("Home", "bootstrap", "static/home"),
		NotFound: views.CreateView("404", "error", "static/404"),
		Error:    views.CreateView("505", "error", "static/505"),
	}
}

type Static struct {
	Home     *views.View
	NotFound *views.View
	Error    *views.View
}

func (c *Static) NewHome(w http.ResponseWriter, r *http.Request) {
	c.Home.Render(w, r, nil)
}
