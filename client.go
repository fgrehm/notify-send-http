package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	summary := "From client"
	body := "Some long and boring body message to be displayed to the user"
	icon := ""
	notificationServer := "http://172.17.42.1:12345"

	flag.StringVar(&icon, "i", icon, "Path to icon")
	flag.String("u", "", "")
	flag.String("a", "", "")
	flag.String("c", "", "")
	flag.String("t", "", "")
	flag.String("h", "", "")
	// TODO: https://github.com/guard/guard/blob/19351271941a3362a47176c6808ddcb4a675e3ad/lib/guard/notifiers/notifysend.rb#L15
	flag.Parse()

	if icon != "" {
		notification := map[string]string{
			"summary": summary,
			"body":    body,
		}
		request, err := newfileUploadRequest(notificationServer, notification, "icon", icon)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		_, err = client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := http.PostForm(notificationServer, url.Values{
			"summary": {summary},
			"body":    {body},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Based on http://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
