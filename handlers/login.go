package handlers

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, r, "login", reCaptchaSiteKey)

	case "POST":
		setToken(w, r)
	}
}
