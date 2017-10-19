package handlers

import (
	"net/http"
)

func events(w http.ResponseWriter, r *http.Request) {
	// TODO(andrei): Finish the events page.
	renderTemplate(w, r, "comingsoon", nil)
}
