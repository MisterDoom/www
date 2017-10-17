package handlers

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates["login"].ExecuteTemplate(w, "layout", reCaptchaSiteKey)

	case "POST":
		email := r.PostFormValue("email")
		// password := r.PostFormValue("password")

		var response string

		response = fmt.Sprintf("Welcome %s", email)
		templates["message"].ExecuteTemplate(w, "layout", messageModel{"Login", response})

	}
}
