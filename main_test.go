package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google-drive-web-app/auth"
	d "google-drive-web-app/drive"
	"google.golang.org/api/drive/v3"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// unit test
func TestHandleLogin(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.HandleLogin)
	handler.ServeHTTP(rr, req)
	print(rr.Header())
	// Check the status code
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}

	// Check the redirect location
	if loc := rr.Header().Get("Location"); loc == "" {
		t.Errorf("handler did not redirect to authentication URL")
	}
}

// this cannot be initiated from here as this callback will be initiated from Google Identity provider
func TestHandleCallback(t *testing.T) {

}

func TestHandleListFiles(t *testing.T) {
	req, err := http.NewRequest("GET", "/listfiles", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(d.HandleListFiles)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "text/html; charset=utf-8")
	}
}

func TestHandleUploadPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/upload", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(d.HandleUploadPage)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleFileUpload(t *testing.T) {
	file, err := os.Open("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/uploadfile", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(d.HandleFileUpload)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleDownload(t *testing.T) {
	service, err := getDriveService()
	if err != nil {
		t.Fatal(err)
	}
	fileName := "testfile.txt"
	query := fmt.Sprintf("name='%s' and mimeType='text/plain' and trashed=false", fileName)
	files, err := service.Files.List().Q(query).Do()
	if err != nil {
		t.Fatal(err)
	}
	fileID := ""
	if len(files.Files) == 0 {
		t.Fatal("no files found")
	}
	fileID = files.Files[0].Id
	req, err := http.NewRequest("GET", fmt.Sprintf("/download?id=%s", fileID), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(d.HandleDownload)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check for Content-Disposition header
	if disp := rr.Header().Get("Content-Disposition"); disp == "" {
		t.Errorf("handler did not set Content-Disposition header")
	}
}

func TestHandleDelete(t *testing.T) {
	service, err := getDriveService()
	if err != nil {
		t.Fatal(err)
	}
	fileName := "testfile.txt"
	query := fmt.Sprintf("name='%s' and mimeType='text/plain' and trashed=false", fileName)
	files, err := service.Files.List().Q(query).Do()
	if err != nil {
		t.Fatal(err)
	}
	fileID := ""
	if len(files.Files) == 0 {
		t.Fatal("no files found")
	}
	fileID = files.Files[0].Id
	req, err := http.NewRequest("GET", fmt.Sprintf("/delete?id=%s", fileID), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(d.HandleDelete)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check for expected output
	expected := "File or folder deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
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

/*
Integration Tests:
The api interaction workflow as:
Auth -> List
Auth -> Upload
Auth -> Download (File information gotten from listfiles page)
Auth -> Delete (File information gotten from listfiles page)
The above tests was run against the application running in localhost and has maintained the above api
workflow to pass the tests.
*/
