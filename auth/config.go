package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     "812135697423-teiluflibfja1gh6bam8a50b79q5rihh.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}
	randomState = "random" // Use a more secure random state in production
)
