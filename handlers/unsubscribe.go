package handlers

import (
	"net/http"

	"github.com/MisterDoom/www/helpers/crypto"
	"github.com/MisterDoom/www/services/databaseService"
)

func unsubscribe(w http.ResponseWriter, r *http.Request) {
	email, err := crypto.Decrypt(r.FormValue("token"))

	if err != nil {
		renderTemplate(w, r, "message", messageModel{"Error", "Token is invalid."})
		return
	}

	if !databaseService.ExistsUser(email) {
		renderTemplate(
			w, r, "message", messageModel{"Error", `Email "` + email + `" is not part of our mailing list.`})
		return
	}

	if err := databaseService.DeleteUser(email); err != nil {
		renderTemplate(w, r, "message", messageModel{"Error", "An unexpected error has occurred. Please try again later."})
		return
	}

	renderTemplate(
		w, r, "message", messageModel{"Unsubscribed", `Email "` + email + `" has been removed from our mailing list.`})
}
