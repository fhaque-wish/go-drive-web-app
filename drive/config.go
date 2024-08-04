package drive

// FileInfo represents the information about a Google Drive file
type FileInfo struct {
	ID           string
	Name         string
	MimeType     string
	ModifiedTime string
}

// TemplateData represents the data to be passed to the HTML template
type TemplateData struct {
	Files []FileInfo
}
