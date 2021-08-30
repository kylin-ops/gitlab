package main

import (
	"bytes"
	"fmt"
	"gitlab"
	"mime/multipart"
	"net/http"
)

var client = gitlab.Client{AccessAddress: "http://192.168.31.220", AccessToken: "JSqmyRboQpNGvWRzxsDx"}

func main() {
	// client.ProjectList()
	// client.ProjectGet("2")
	client.ProjectCreate("test")

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fmt.Println(w.FormDataContentType())
	http.PostForm()
}
