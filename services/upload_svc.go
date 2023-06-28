package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadToBucket(f2 multipart.File, fname string) string {
	// Create a new buffer to store the form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a new form file field
	fileField, err := writer.CreateFormFile("file", fname)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Copy the file data to the form file field
	_, err = io.Copy(fileField, f2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	token := "f9403fc5f537b4ab332d"
	_ = writer.WriteField("token", token)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	request, err := http.NewRequest("POST", "http://localhost:25478/upload", body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer response.Body.Close()

	// Print the response
	fmt.Println(response.Status)
	fmt.Println(response)

	if response.StatusCode == 200 {
		url := "http://localhost:25478/files/" + fname + "?token=" + token
		return url
	}
	return ""
}
