package drive

import (
	"fmt"
	"html/template"
	"net/http"
)

// HandleListFiles renders the list of files in Google Drive in HTML format
func HandleListFiles(w http.ResponseWriter, r *http.Request) {
	driveService, err := getDriveService()
	if err != nil {
		err = fmt.Errorf("Failed to create drive client:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	renderListTemplate(&w, &fileInfos)
}

func renderListTemplate(w *http.ResponseWriter, f *[]FileInfo) {
	tmpl, err := template.ParseFiles("templates/files.html")
	if err != nil {
		http.Error(*w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	//serving data to html
	data := TemplateData{Files: *f}
	err = tmpl.Execute(*w, data)
	if err != nil {
		http.Error(*w, "Failed to render template", http.StatusInternalServerError)
		return
	}

}
