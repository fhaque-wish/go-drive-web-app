package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     getEnv("OAUTHCLIENTID", ""),
		ClientSecret: getEnv("OAUTHCLIENTSECRET", ""),
		RedirectURL:  getEnv("REDIRECTURL", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}
	randomState = "random" // Use a more secure random state in production
)
