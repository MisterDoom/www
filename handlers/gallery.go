package handlers

import (
	"net/http"

	"github.com/MisterDoom/www/services/galleryService"
)

func gallery(w http.ResponseWriter, r *http.Request) {
	albums := galleryService.GetAlbums()
	templates["gallery"].ExecuteTemplate(w, "layout", &albums)
}
