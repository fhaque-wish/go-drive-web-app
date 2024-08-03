package drive

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"os"
)

func getToken() (*oauth2.Token, error) {
	f, err := os.Open("token.json")
	if err != nil {
		log.Error("error opening token.json")
		return nil, err
	}
	defer f.Close()
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}

/*func initClient() *client {
	newClient := client{}
	token, err := getToken()
	if err != nil {
		return &newClient
	}
	newClient.token = token.AccessToken
	newClient.c = auth.GoogleOauthConfig.Client(context.Background(), token)
	return &newClient
}*/
