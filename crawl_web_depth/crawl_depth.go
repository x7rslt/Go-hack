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
	StartUrl   *url.URL
	timeout    = time.Duration(8 * time.Second)
)

func Test(path string) {
	fmt.Println("Test:", StartUrl, "Test2:", StartUrl.Scheme)
	fmt.Println(path)
}
func CrawlUrl(path string) {
	fmt.Println("Crawl url:", StartUrl.String(), StartUrl.Scheme)
	var targetUrl url.URL
	targetUrl.Scheme = StartUrl.Scheme
	targetUrl.Host = StartUrl.Host
	targetUrl.Path = path
	httpClient := http.Client{Timeout: timeout}
	response, err := httpClient.Get(targetUrl.String())
	if err != nil {
		log.Println("Response error :", err, "Http get url :", targetUrl.String())
	}
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
	if tempUrl.Host != "" && tempUrl.Host != StartUrl.Host {
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
	StartUrl, err := url.Parse("https://www.eastmoney.com/")
	if err != nil {
		log.Fatal("Error parsing URL.", err)
	}

	fmt.Println("Starturl:", StartUrl.Scheme, StartUrl.Host, StartUrl.Path)
	fmt.Printf("Starturl type %T", StartUrl)
	//CrawlUrl(StartUrl.Path)
	fmt.Println(foundPaths)
	log.Println("Found paths number:", len(foundPaths))
}
