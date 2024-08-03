package drive

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// HandleDownload handles file download requests
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "File ID not specified", http.StatusBadRequest)
		return
	}

	driveService, err := getDriveService()
	if err != nil {
		err = fmt.Errorf("Failed to create drive client:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get file metadata
	file, err := driveService.Files.Get(fileID).Fields("name").Do()
	if err != nil {
		http.Error(w, "Failed to get file metadata", http.StatusInternalServerError)
		return
	}

	// Download the file
	response, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		err = fmt.Errorf("Failed to download file:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Create a local file
	localFilePath := filepath.Join(os.TempDir(), file.Name)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		http.Error(w, "Failed to create local file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Copy the content to the local file
	_, err = io.Copy(localFile, response.Body)
	if err != nil {
		http.Error(w, "Failed to save file locally", http.StatusInternalServerError)
		return
	}

	// Serve the file to the user
	localFile.Seek(0, 0)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, localFile)
}
