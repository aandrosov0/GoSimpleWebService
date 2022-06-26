package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	var err error

	switch request.Method {
	case "GET":
		err = handleGet(writer, request)
	case "POST":
		err = handlePost(writer, request)
	case "PUT":
		err = handlePut(writer, request)
	case "DELETE":
		err = handleDelete(writer, request)
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func handleGet(writer http.ResponseWriter, request *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(request.URL.Path))

	if err != nil {
		return
	}

	post, err := retrieve(id)

	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t\t")

	if err != nil {
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(output)
	return
}

func handlePost(writer http.ResponseWriter, request *http.Request) (err error) {
	post := Post{}
	len := request.ContentLength
	body := make([]byte, len)

	request.Body.Read(body)
	json.Unmarshal(body, &post)

	if err = post.create(); err != nil {
		return
	}

	writer.WriteHeader(200)
	return
}

func handlePut(writer http.ResponseWriter, request *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(request.URL.Path))

	if err != nil {
		return
	}

	post, err := retrieve(id)

	if err != nil {
		return
	}

	len := request.ContentLength
	body := make([]byte, len)

	request.Body.Read(body)
	json.Unmarshal(body, &post)

	if err = post.update(); err != nil {
		return
	}

	writer.WriteHeader(200)
	return
}

func handleDelete(writer http.ResponseWriter, request *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(request.URL.Path))

	if err != nil {
		return
	}

	post, err := retrieve(id)

	if err != nil {
		return
	}

	if err = post.delete(); err != nil {
		return
	}

	writer.WriteHeader(200)
	return
}

func main() {
	http.HandleFunc("/post/", handleRequest)
	http.ListenAndServe(":8080", nil)
}
