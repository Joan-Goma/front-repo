package liro

import (
	"net/http"

	"neft.web/views"
)

type Liro struct {
	LiroMain     *views.View
	LiroUserList *views.View
}

func NewLiro() *Liro {
	views.Templates["LiroMain"] = views.Template{
		Layout: "dashboard",
		File:   "liro/index",
		View:   views.NewView("dashboard", "liro/index"),
	}
	views.Templates["LiroUserList"] = views.Template{
		Layout: "dashboard",
		File:   "liro/listUsers",
		View:   views.NewView("dashboard", "liro/listUsers"),
	}
	return &Liro{
		LiroMain:     views.Templates["LiroMain"].View,
		LiroUserList: views.Templates["LiroUserList"].View,
	}
}

// New GET /liro
func (u *Liro) New(w http.ResponseWriter, r *http.Request) {
	u.LiroMain.Render(w, r, nil)
}

// New GET /liro/users
func (u *Liro) UsersList(w http.ResponseWriter, r *http.Request) {
	u.LiroUserList.Render(w, r, nil)
}
