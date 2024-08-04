package drive

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// HandleUploadPage renders the upload file page
func HandleUploadPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/upload.html")
}

// HandleFileUpload processes the file upload and uploads it to Google Drive
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the target folder from the form data
	targetFolder := r.FormValue("folder")

	// Save the file locally (optional, for further processing if needed)
	localFilePath := filepath.Join(os.TempDir(), handler.Filename)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		http.Error(w, "Error creating local file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, file)
	if err != nil {
		http.Error(w, "Error saving local file", http.StatusInternalServerError)
		return
	}

	driveService, err := getDriveService()
	if err != nil {
		err = fmt.Errorf("Failed to create drive client: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the target folder exists, if not, create it
	folderID, err := getOrCreateFolder(driveService, targetFolder)
	if err != nil {
		http.Error(w, "Failed to get or create folder", http.StatusInternalServerError)
		return
	}

	// Perform resumable upload
	driveFile := &drive.File{
		Name: handler.Filename,
		//Parents: []string{folderID},
	}
	//upload into a specific folder
	if folderID != "" {
		driveFile.Parents = []string{folderID}
	}

	localFile, err = os.Open(localFilePath)
	if err != nil {
		http.Error(w, "Error opening local file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	//resumble upload facilitates larger file upload
	resumableUpload := driveService.Files.Create(driveFile).ResumableMedia(context.Background(), localFile, handler.Size, handler.Header.Get("Content-Type"))
	uploadedFile, err := resumableUpload.Do()
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File uploaded successfully to folder '%s': %s", targetFolder, uploadedFile.Name)
}
