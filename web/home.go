package web

import (
	"net/http"
)

// HandleHome serves home.html
func HandleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
}
