package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
	"fmt"
)

const baseURL = "https://www.economist.com"
const weeklyEditionURL = baseURL + "/weeklyedition"

type Section struct {
	title string
	urls []string
}

type Article struct {
	supertitle string
	title string
	subtitle string
	body string
}

func main() {
	doc := getDocument(weeklyEditionURL)

	markdown := buildMarkdown(doc)
	fmt.Println(markdown)
}

func getDocument(url string) goquery.Document {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// return response.Body
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return *doc
}

func headline(doc goquery.Document) string {
	return doc.Find(".weekly-edition-header__headline").Text()
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

func getArticle(url string) Article {
	doc := getDocument(url)

	bodyBuilder := strings.Builder{}

	supertitle := doc.Find(".css-1mdqtqm.e11bnqe00").Text()
	title := doc.Find(".css-6p5r16.e11fb5bd0").Text()
	subtitle := doc.Find(".css-qm1vln.elmr55g0").Text()

	doc.Find(".css-1qkdneh.ee5d8yd2").Children().Filter(":not(style)").Each(func(_ int, s *goquery.Selection) {
		// I hate pseudoselectors
		if isRealContent(s) {
			if s.Is("h2") {
				bodyBuilder.WriteString("#### " + s.Text() + "\n\n")
			} else if s.Is("p") {
				bodyBuilder.WriteString(s.Text() + "\n\n")
			} else if s.Is("div") && (s.HasClass("css-0") || s.HasClass("css-1bzwzkr")) {
				imageURL := s.Find("img").AttrOr("src", "")
				imageCaption := s.Find(".css-6pzdis.edweubf1").Text()
				bodyBuilder.WriteString("![" + imageCaption + "](" + imageURL + ")\n\n")
			}
		}
	})

	return Article{supertitle, title, subtitle, bodyBuilder.String()}
}

func isRealContent(s *goquery.Selection) bool {
	return !s.HasClass("css-1lm38nn") && !s.HasClass("adComponent_advert__V79Pp")
}

func buildMarkdown(doc goquery.Document) string {
	headline := headline(doc)
	worldThisWeek := worldThisWeek(doc)
	leaders := leaders(doc)
	sections := sections(doc)

	builder := strings.Builder{}

	builder.WriteString("---\ntitle: \"" + headline + "\"\n---\n\n")

	builder.WriteString("# The world this week\n\n")
	for _, url := range worldThisWeek {
		article := getArticle(url)
		builder.WriteString("## " + article.title + "\n\n")
		builder.WriteString(article.body)
	}

	builder.WriteString("# Leaders\n\n")
	for _, url := range leaders {
		article := getArticle(url)
		builder.WriteString("## " + article.title + "\n\n")
		builder.WriteString(article.body)
	}

	for _, section := range sections {
		builder.WriteString("# " + section.title + "\n\n")
		for _, url := range section.urls {
			article := getArticle(url)
			builder.WriteString("## " + article.title + "\n\n")
			builder.WriteString(article.body)
		}
	}

	return builder.String()
}

