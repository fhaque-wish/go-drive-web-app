package main

import (
	"fmt"
	"google-drive-web-app/auth"
	"google-drive-web-app/drive"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", auth.HandleMain)
	http.HandleFunc("/home", auth.HandleHome)
	http.HandleFunc("/login", auth.HandleLogin)
	http.HandleFunc("/callback", auth.HandleCallback)
	http.HandleFunc("/listfiles", drive.HandleListFiles)
	http.HandleFunc("/upload", drive.HandleUploadPage)
	http.HandleFunc("/uploadfile", drive.HandleFileUpload)
	http.HandleFunc("/download", drive.HandleDownload)
	http.HandleFunc("/delete", drive.HandleDelete)

	// Serve static files (HTML templates)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
