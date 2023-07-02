package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadToBucket(f2 multipart.File, fname string) string {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	fileField, err := writer.CreateFormFile("file", fname)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	_, err = io.Copy(fileField, f2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	token := os.Getenv("BUCKET_TOKEN")
	bucketUrl := os.Getenv("BUCKET_URL")
	_ = writer.WriteField("token", token)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	request, err := http.NewRequest("POST", bucketUrl+"/upload", body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		url := bucketUrl + "/files/" + fname + "?token=" + token
		return url
	}
	return ""
}
