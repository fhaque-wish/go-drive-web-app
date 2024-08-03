package drive

import (
	"fmt"
	"net/http"
)

// HandleDelete handles file deletion requests
func HandleDelete(w http.ResponseWriter, r *http.Request) {
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

	// Delete the file
	err = driveService.Files.Delete(fileID).Do()
	if err != nil {
		http.Error(w, "Failed to delete file or folder", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File or folder deleted successfully")
}
