package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func hello(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := TodoPageData{}
	tmpl.Execute(w, data)
}

func robots(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "User-agent: *")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/robots.txt", robots).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":80", nil)
}
