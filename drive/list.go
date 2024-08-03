package drive

import (
	"context"
	"google-drive-web-app/auth"
	"google.golang.org/api/drive/v3"
	"html/template"
	"net/http"
)

// HandleListFiles renders the list of files in Google Drive in HTML format
func HandleListFiles(w http.ResponseWriter, r *http.Request) {
	token, err := getToken()
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	client := auth.GoogleOauthConfig.Client(context.Background(), token)
	driveService, err := drive.New(client)
	if err != nil {
		http.Error(w, "Failed to create drive client", http.StatusInternalServerError)
		return
	}

	files, err := driveService.Files.List().Fields("files(id, name, mimeType, modifiedTime)").Do()
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	var fileInfos []FileInfo
	for _, file := range files.Files {
		fileInfos = append(fileInfos, FileInfo{
			ID:           file.Id,
			Name:         file.Name,
			MimeType:     file.MimeType,
			ModifiedTime: file.ModifiedTime,
		})
	}

	tmpl, err := template.ParseFiles("templates/files.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := TemplateData{Files: fileInfos}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
