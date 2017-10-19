package handlers

import (
	"context"
	"net/http"
	"time"
)

func logout(w http.ResponseWriter, r *http.Request) {
	expireCookie := time.Now().Add(-time.Hour * 24)
	cookie := http.Cookie{Name: "Auth", Value: "expired", Expires: expireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)

	user := &User{
		LoggedIn: false,
	}

	ctx := context.WithValue(r.Context(), "user", user)

	renderTemplate(w, r.WithContext(ctx), "message", &messageModel{Title: "Logout OK", Message: "You've logged out successfully"})
}
