package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"neft.web/errorController"

	"neft.web/views"
)

func NewUsers() *Users {
	return &Users{
		NewView:      views.NewView("dashboard", "users/register"),
		LoginView:    views.NewView("dashboard", "users/login"),
		ForgotPwView: views.NewView("dashboard", "users/forgot_pw"),
		ResetPwView:  views.NewView("dashboard", "users/reset_pw"),
	}
}

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
}

// New GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New GET /login
func (u *Users) LoginNew(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

// New POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "neftAuth",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Create Process the signup form
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	// TODO
	var vd views.Data

	var form SignupForm
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.NewView.Render(w, r, &vd)
		errorController.HandleError(err.Error(), "Parsing Register Form")
		return
	}

	vd.Yield = &form
	personJSON, err := processFormToAPI(r)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.NewView.Render(w, r, &vd)
	}

  req, err := http.NewRequest(http.MethodPut, "https://test.joan-goma.repl.co/v1/auth", bytes.NewBuffer(personJSON))
    if err != nil {
        // handle error
    }
  
    req.Header.Set("Content-Type", "application/json")
	
   client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
    vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		errorController.ErrorLogger.Println("Error sending request login")
		u.NewView.Render(w, r, &vd)
		return
    }
    defer resp.Body.Close()

	if resp.StatusCode != 201 {
		answer, err := readAPIAnswer(resp)
		if err != nil {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: template.HTML(views.AlertMsgGeneric),
			}
			u.NewView.Render(w, r, &vd)
			return
		}
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: template.HTML(answer.ErrorField),
		}
		u.NewView.Render(w, r, &vd)
		return
	}

	if err := u.signIn(w, resp.Header.Get("neftAuth")); err != nil {
		vd.Alert = &views.Alert{
			Level:   vd.Alert.Level,
			Message: template.HTML(views.AlertMsgGeneric),
		}
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Sign In error generating Remember token"
		errorController.WD.SendErrorWHWeb()
		return
	}
	http.Redirect(w, r, "/liro", http.StatusFound)
}

type LoginForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var loginData LoginForm
	if err := ParseForm(r, &loginData); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.LoginView.Render(w, r, &vd)
		return
	}
	vd.Yield = loginData
	personJSON, err := processFormToAPI(r)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.LoginView.Render(w, r, &vd)
	}

	resp, err := http.Post("https://test.joan-goma.repl.co/api/login", "application/json", bytes.NewBuffer(personJSON))
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		errorController.ErrorLogger.Println("Error sending request login")
		u.LoginView.Render(w, r, &vd)
		return
	}

	if resp.StatusCode != 200 {
		answer, err := readAPIAnswer(resp)
		if err != nil {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: template.HTML(views.AlertMsgGeneric),
			}
			u.LoginView.Render(w, r, &vd)
			return
		}
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: template.HTML(answer.ErrorField),
		}
		u.LoginView.Render(w, r, &vd)
		return
	}

	if err := u.signIn(w, resp.Header.Get("neftAuth")); err != nil {
		vd.Alert = &views.Alert{
			Level:   vd.Alert.Level,
			Message: template.HTML(views.AlertMsgGeneric),
		}
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Sign In error generating Remember token"
		errorController.WD.SendErrorWHWeb()
		return
	}
	http.Redirect(w, r, "/liro", http.StatusFound)
}

type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

//POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	// var vd views.Data
	// var form ResetPwForm
	// vd.Yield = &form
	// if err := ParseForm(r, &form); err != nil {
	// 	vd.Alert = &views.Alert{
	// 		Level:   views.AlertLvlError,
	// 		Message: template.HTML(err.Error()),
	// 	}
	// 	u.ForgotPwView.Render(w, r, &vd)
	// 	return
	// }

	// user, err := u.us.ByEmail(form.Email)
	// if err != nil {
	// 	vd.Alert = &views.Alert{
	// 		Level:   views.AlertLvlError,
	// 		Message: template.HTML(models.ErrEmailNotExist.Error()),
	// 	}
	// 	u.ForgotPwView.Render(w, r, &vd)
	// 	return
	// }
	// token, err := u.us.InitiateReset(user.ID)
	// if err != nil {
	// 	errorController.WD.Content = err.Error()
	// 	errorController.WD.Site = "Error trying to initiate reset"
	// 	errorController.WD.SendErrorWHWeb()
	// 	vd.Alert = &views.Alert{
	// 		Level:   views.AlertLvlError,
	// 		Message: views.AlertMsgGeneric,
	// 	}
	// 	u.ForgotPwView.Render(w, r, &vd)
	// 	return
	// }
	// _ = token
	// vd.Alert = &views.Alert{
	// 	Level:   views.AlertLvlSuccess,
	// 	Message: "Congratulations, you will receive instructions to reset your password via email",
	// }
	// u.ForgotPwView.Render(w, r, &vd)
	http.Redirect(w, r, "/liro", http.StatusFound)
}

//GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	u.ResetPwView.Render(w, r, &vd)
}

//POST /RESET

func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	// user, err := u.us.CompleteReset(form.Token, form.Password)
	// if err != nil {
	// 	if err == models.ErrSamePasswordReset {
	// 		vd.Alert = &views.Alert{
	// 			Level:   views.AlertLvlError,
	// 			Message: template.HTML(models.ErrSamePasswordReset.Error()),
	// 		}
	// 	} else {
	// 		errorController.WD.Content = err.Error()
	// 		errorController.WD.Site = "Error completing the password recovering"
	// 		errorController.WD.SendErrorWHWeb()
	// 		vd.Alert = &views.Alert{
	// 			Level:   views.AlertLvlError,
	// 			Message: views.AlertMsgGeneric,
	// 		}
	// 	}
	// 	u.ResetPwView.Render(w, r, &vd)
	// 	return
	// }
	//u.signIn(w, user)
	http.Redirect(w, r, "/liro", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, token string) error {
	cookie := http.Cookie{
		Name:     "neftAuth",
		Value:    token,
		HttpOnly: false,
		Expires:  time.Now().Local().Add(12 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	return nil
}
