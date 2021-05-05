package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var documentExtensions = []string{"doc", "pdf", "ppt", "txt", "docx", "csv", "xls", "zip", "gz", "tar"}

func TestFindDocument(t *testing.T) {
	target := "https://www.baidu.com"
	fmt.Println("Start search document from ", target)
	resp, err := http.Get(target)
	if err != nil {
		fmt.Println("http Get err:", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("goquery err:", err)
	}
	documtenturls := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && linkContainDocument(href) {
			documtenturls = append(documtenturls, href)
		}
	})
	if len(documtenturls) != 0 {
		fmt.Println("Found Document:", documtenturls)
	} else {
		fmt.Println("No found.")
	}

}

func linkContainDocument(url string) bool {
	urlSplit := strings.Split(url, ".")
	if len(urlSplit) < 2 {
		return false
	}

	for _, ext := range documentExtensions {
		if ext == urlSplit[len(urlSplit)-1] {
			return true
		}
	}
	return false
}
