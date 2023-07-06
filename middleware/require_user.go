package middleware

import (
	"errors"
	"net/http"
	"time"

	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Sabadell0310JED")

type JWTClaim struct {
	RemmemberHash string `json:"remmemberHash"`
	RoleID        int    `json:"role"`
	jwt.StandardClaims
}

type User struct {
}

func (mw *User) Apply(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	})
}

type RequireUser struct {
	User
}

func (mw *RequireUser) Apply(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nefToken, err := r.Cookie("neftAuth")

		if err != nil {
			next(w, r)
			return
		}

		if nefToken.Value == "" {
			next(w, r)
			return
		}

		err = validateToken(nefToken.Value)

		if err != nil {
			fmt.Println(err)
			c := &http.Cookie{
				Name:     "neftAuth",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, c)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/liro", http.StatusFound)
	})
}
func (mw *RequireUser) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nefToken, err := r.Cookie("neftAuth")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if nefToken.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		err = validateToken(nefToken.Value)

		if err != nil {
			fmt.Println(err)
			c := &http.Cookie{
				Name:     "neftAuth",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, c)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}

func validateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt <= time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
