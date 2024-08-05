# go-drive-web-app

`go-drive-web-app` is a web application that interacts with a user's Google Drive. The application performs the following tasks:

1. Authenticate the user using OAuth 2.0.
2. List files in the user’s Google Drive.
3. Upload a file to the user’s Google Drive.
4. Download a file from the user’s Google Drive.
5. Delete a file from the user’s Google Drive.

## Development Setup

1. Configure a Google Cloud project by following the instructions [here](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).
2. In the project, enable the Google Drive API from API & Services.
3. Set up an OAuth 2.0 Client ID from Credentials.
4. Set up the OAuth consent screen.
5. Allow Google Drive API scope access to the OAuth client.
6. Optionally, specify allowed test users.
7. Once the OAuth client is set up, it will provide a Client ID and Client Secret.
8. Create a `.env` file in this project's root folder and set the following values:

```text
OAUTHCLIENTID=<OAuth client ID>
OAUTHCLIENTSECRET=<OAuth client secret>
REDIRECTURL=<OAuth redirect URL>
```

## Run the Application

Run the following command from the project's root folder:

```go
go run main.go
```
The application will start at  http://localhost:8080

## Test Setup
The `main_test.go` file contains the test cases. To run tests:
1. Ensure the application is running.
2. Ensure a test user is authenticated with the application and has access to a test google drive.
3. Create a testfile.txt file in the project's root folder. This file is used to test upload, download, and delete functionalities. If there is a file name conflict in the test google drive, rename the `testfile.txt` to a unique name and update the references in the `main_test.go` file.
4. Run the following command in a separate terminal:
```go
go test -v
```

