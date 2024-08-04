package drive

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google-drive-web-app/auth"
	"google.golang.org/api/drive/v3"
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

func getDriveService() (*drive.Service, error) {
	driveService := drive.Service{}
	token, err := getToken()
	if err != nil {
		log.Error("error getting token")
		return &driveService, err
	}

	client := auth.GoogleOauthConfig.Client(context.Background(), token)
	return drive.New(client)
}

// getOrCreateFolder checks if a folder exists and returns its ID, otherwise creates it
func getOrCreateFolder(service *drive.Service, folderName string) (string, error) {
	if folderName == "" {
		return "", nil
	}
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", folderName)
	files, err := service.Files.List().Q(query).Do()
	if err != nil {
		return "", err
	}

	if len(files.Files) > 0 {
		return files.Files[0].Id, nil
	}

	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return "", err
	}

	return createdFolder.Id, nil
}
