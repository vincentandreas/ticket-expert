package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func UploadToBucket(f2 multipart.File, fname string) string {
	body := bytes.NewBuffer(nil)

	_, err := io.Copy(body, f2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	bucketUrl := os.Getenv("BUCKET_URL")
	bucketName := os.Getenv("BUCKET_NAME")

	fileName := uuid.New().String()


	fullUrl := bucketUrl + "/" + bucketName + "/" + fileName

	request, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(body.Bytes()))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return fullUrl
	}
	return ""
}
