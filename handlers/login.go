package handlers

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates["login"].ExecuteTemplate(w, "layout", reCaptchaSiteKey)
	}

}
