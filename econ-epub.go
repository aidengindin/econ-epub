package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://www.economist.com"
const weeklyEditionURL = baseURL + "/weeklyedition"

type Section struct {
	title string
	urls []string
}

func main() {
	body := makeRequest()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	// headline := headline(*doc)
	// worldThisWeek := worldThisWeek(*doc)
	// leaders := leaders(*doc)
	// sections := sections(*doc)
}

func makeRequest() io.Reader {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(weeklyEditionURL)
	if err != nil {
		log.Fatal(err)
	}
	// defer response.Body.Close()
	return response.Body
}

func headline(doc goquery.Document) string {
	var headline string
	doc.Find(".weekly-edition-header__headline").Each(func(i int, s *goquery.Selection) {
		headline = s.Text()
	})
	return headline
}

func worldThisWeek(doc goquery.Document) []string {
	var articles []string
	doc.Find(".weekly-edition-wtw__item").Each(func(i int, s *goquery.Selection) {
		articles = append(articles, baseURL + s.Find("a").AttrOr("href", ""))
	})
	return articles
}

func leaders(doc goquery.Document) []string {
	var articles []string
	doc.Find(".teaser-weekly-edition--leaders").Find(".css-qakuwj.e1rr6cni0").Each(func(_ int, s *goquery.Selection) {
		articles = append(articles, baseURL + s.Find("a").AttrOr("href", ""))
	})
	return articles
}

func sections(doc goquery.Document) []Section {
	var sections []Section
	doc.Find(".layout-weekly-edition-section").Each(func(_ int, s *goquery.Selection) {
		var urls []string;
		title := s.Find(".ds-section-headline").Text()
		s.Find(".css-wl43hz.e1rr6cni0").Each(func(_ int, t *goquery.Selection) {
			urls = append(urls, baseURL + t.Find("a").AttrOr("href", ""))
		})
		sections = append(sections, Section{title, urls})
	})
	return sections
}

