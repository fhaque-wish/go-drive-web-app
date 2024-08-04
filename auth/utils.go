package auth

import (
	"encoding/json"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
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

// getEnv reads an environment variable and provides a default value if not set
func getEnv(key, defaultValue string) string {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
}
