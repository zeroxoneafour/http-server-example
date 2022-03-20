package main

import (
	"fmt"
	"log"
	"mime"
	"os"
	"strings"

	http "github.com/zeroxoneafour/http-server"
)

func readFile(path string) (status http.Status, size int64, filename, mimetype, content string) {
	filename = "." + path
	if fileinfo, err := os.Stat(filename); err == nil {
		if fileinfo.IsDir() {
			if filename[len(filename)-1] != '/' {
				filename += "/"
			}
			filename += "index.html"
		}
	} else {
		status = 404
		return
	}
	file, err := os.Open(filename)
	defer file.Close()
	fileinfo, _ := os.Stat(filename)
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		status = 403
		return
	}
	content = string(contentBytes)
	size = fileinfo.Size()
	mimetype = mime.TypeByExtension(filename[strings.LastIndex(filename, "."):]) // get mimetype of extension of filename
	status = 200
	return
}

func getHandler(client *http.HTTPClient) http.Status {
	status, size, filename, mimetype, content := readFile(client.Req.GetPath())
	client.Res.Content = content
	client.Res.Headers["Content-Length"] = fmt.Sprint(size)
	client.Res.Headers["Content-Type"] = mimetype
	log.Println("GET request of", filename[1:], "returned", fmt.Sprint(status))
	return status
}

func headHandler(client *http.HTTPClient) http.Status {
	status, size, filename, mimetype, _ := readFile(client.Req.GetPath())
	client.Res.Headers["Content-Length"] = fmt.Sprint(size)
	client.Res.Headers["Content-Type"] = mimetype
	log.Println("HEAD request of", filename[1:], "returned", fmt.Sprint(status))
	return status
}

func main() {
	host := "localhost"
	port := "8000"
	if val, present := os.LookupEnv("HOST"); present {
		host = val
	}
	if val, present := os.LookupEnv("PORT"); present {
		port = val
	}

	server := http.New(host, port)
	server.SetHandler(http.GET, getHandler)
	server.SetHandler(http.HEAD, headHandler)
	log.Print("Starting server on ", host, ":", port)
	server.Run()
}
