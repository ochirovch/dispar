package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Settings struct {
	ProjectName         string
	PaginationType      string
	PaginationURL       string
	PaginationStartPage int
	PaginationEndPage   int
	CurrentPage         int
	LinkXPath           string
	ContentXPath        string
	DataXPath           []string // Data in Content
}

func hello(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := Settings{}
	tmpl.Execute(w, data)
}

func robots(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "User-agent: *")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/robots.txt", robots).Methods("GET")

	http.ListenAndServe(":80", nil)
}
