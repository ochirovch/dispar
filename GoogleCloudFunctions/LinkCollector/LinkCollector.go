package LinkCollector

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"
)

const (
	begin = 1
	end   = 8
)

func LinkCollector(w http.ResponseWriter, req *http.Request) {
	c := colly.NewCollector()

	// add data to pub sub channel page project;url
	c.OnHTML(".news-item__body > a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := begin; i < end; i++ {
		fmt.Println(i)
		c.Visit("https://www.chita.ru/news/?pg=" + strconv.Itoa(i))
	}
}
