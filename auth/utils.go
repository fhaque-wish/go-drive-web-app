package auth

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"os"
)

func saveToken(token *oauth2.Token) error {
	f, err := os.Create("token.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
