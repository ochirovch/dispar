package LinkCollector

import (
	"net/http"

	"github.com/gocolly/colly/v2"
)

func LinkCollector(w http.ResponseWriter, req *http.Request) {
	c := colly.NewCollector()

	// add data to pub sub channel page project;url
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.Visit("http://go-colly.org/")
}
