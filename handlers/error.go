package handlers

import "net/http"

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	if status == http.StatusNotFound {
		renderTemplate(w, r, "message", messageModel{"Error", "Page not found."})
	} else {
		renderTemplate(w, r, "message", messageModel{"Error", "An error has occurred."})
	}
}
