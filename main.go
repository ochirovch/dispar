package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	colly "github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
)

const (
	begin = 1
	end   = 100
)

type Settings struct {
	ProjectName         string
	PaginationType      string
	PaginationURL       string
	PaginationStartPage int
	PaginationEndPage   int
	CurrentPage         int
	LinkSelector        string
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

func LinkCollector(w http.ResponseWriter, req *http.Request) {

	c := colly.NewCollector()

	// add data to pub sub channel page project;url

	c.OnHTML("div.news-item__title > a", func(e *colly.HTMLElement) {
		fmt.Println(e.Request.URL.Host + e.Attr("href"))
		//		e.Request.Visit(e.Attr("href"))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := begin; i < end; i++ {
		fmt.Println(i)
		c.Visit("https://www.chita.ru/news/?pg=" + strconv.Itoa(i))
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/robots.txt", robots).Methods("GET")
	router.HandleFunc("/linkcollector", LinkCollector).Methods("GET")

	http.ListenAndServe(":8010", router)
}
