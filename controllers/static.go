package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"neft.web/errorController"

	"neft.web/views"
)

func NewStatic() *Static {
	views.Templates["Home"] = views.Template{
		Layout: "bootstrap",
		File:   "static/home",
		View:   views.NewView("bootstrap", "static/home"),
	}
	views.Templates["404"] = views.Template{
		Layout: "error",
		File:   "static/404",
		View:   views.NewView("error", "static/404"),
	}
	views.Templates["505"] = views.Template{
		Layout: "error",
		File:   "static/505",
		View:   views.NewView("error", "static/505"),
	}
	return &Static{
		Home:     views.Templates["Home"].View,
		NotFound: views.Templates["404"].View,
		Error:    views.Templates["505"].View,
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

func NewContact() *Contact {
	return &Contact{
		HomeView:  views.NewView("bootstrap", "static/home"),
		LoginView: views.NewView("dashboard", "users/login"),
	}
}

type Contact struct {
	HomeView  *views.View
	LoginView *views.View
}

type ContactForm struct {
	Name    string `schema:"name"`
	Email   string `schema:"email"`
	Subject string `schema:"subject"`
	Message string `schema:"message"`
}

// Create Process the contact form
// POST /
func (c *Contact) ContactForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	var vd views.Data
	var form ContactForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Error parsing the contact form"
		errorController.WD.SendErrorWHWeb()
		return
	}
	form.Email = strings.ToLower(form.Email)
	form.Email = strings.TrimSpace(form.Email)

	if !emailRegex.MatchString(form.Email) {
		vd.SetAlert("Contact mail is not valid!")
		c.HomeView.Render(w, r, &vd)
		return
	}

	form = ContactForm{}
	c.HomeView.Flush(w, r, &vd)
}
