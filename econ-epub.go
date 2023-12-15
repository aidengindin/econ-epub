package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get("https://www.economist.com/weeklyedition")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".weekly-edition-header__headline").Each(func(i int, s *goquery.Selection) {
		log.Println(s.Text())
	})
}

