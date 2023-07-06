package liro

import (
	"net/http"

	"neft.web/views"
)

type Liro struct {
	Main     *views.View
	UserList *views.View
}

func NewLiro() *Liro {
	return &Liro{
		Main:     views.NewView("dashboard", "liro/index"),
		UserList: views.NewView("dashboard", "liro/listUsers"),
	}
}

// New GET /liro
func (u *Liro) New(w http.ResponseWriter, r *http.Request) {
	u.Main.Render(w, r, nil)
}

// New GET /liro/users
func (u *Liro) UsersList(w http.ResponseWriter, r *http.Request) {
	u.UserList.Render(w, r, nil)
}
