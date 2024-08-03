package auth

import (
	"context"
	"fmt"
	"net/http"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><body>
        <a href="/login">Google Log In</a>
        </body></html>`)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := GoogleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	//client := googleOauthConfig.Client(context.Background(), token)
	// Save the token to a file or session for later use
	err = saveToken(token)
	if err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "authentication successful")
	//http.Redirect(w, r, "/listfiles", http.StatusTemporaryRedirect)
}
