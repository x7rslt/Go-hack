package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	foundPaths []string
	startUrl   url.URL
	timeout    = time.Duration(8 * time.Second)
)

func Test(path string) {
	fmt.Println("Test:", startUrl, "Test2:", startUrl.Scheme)
	fmt.Println(path)
}
func CrawlUrl(path string) {
	fmt.Println(startUrl, startUrl.Scheme)
	var targetUrl url.URL
	targetUrl.Scheme = startUrl.Scheme
	targetUrl.Host = startUrl.Host
	targetUrl.Path = path
	httpClient := http.Client{Timeout: timeout}
	response, err := httpClient.Get(targetUrl.String())
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		parseUrl, err := url.Parse(href)
		if err != nil {
			return
		}
		if UrlIsInScope(parseUrl) {
			foundPaths = append(foundPaths, parseUrl.Path)
			log.Println("Found new path to crawl:", parseUrl.String())
			CrawlUrl(parseUrl.Path)
		}
	})

}

func UrlIsInScope(tempUrl *url.URL) bool {
	if tempUrl.Host != "" && tempUrl.Host != startUrl.Host {
		return false
	}
	if tempUrl.Path == "" {
		return false
	}
	for _, existsPath := range foundPaths {
		if existsPath == tempUrl.Path {
			return false
		}
	}
	return true
}

func main() {
	foundPaths = make([]string, 0)
	startUrl, err := url.Parse("https://www.eastmoney.com/")
	if err != nil {
		log.Fatal("Error parsing URL.", err)
	}
	log.Println("Crawling:", startUrl.String())
	fmt.Println(startUrl.Scheme, startUrl.Host, startUrl.Path)
	CrawlUrl(startUrl.Path)
	fmt.Println(foundPaths)
	log.Println("Found paths number:", len(foundPaths))
}
