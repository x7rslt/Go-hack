package test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func checkfile(baseurl, path string, done chan bool) {
	targetUrl, err := url.Parse(baseurl)
	if err != nil {
		log.Println("Error parse base Url.", err)
	}
	targetUrl.Path = path
	res, err := http.Head(targetUrl.String())
	if err != nil {
		log.Println("Error fetching", targetUrl.String())
	}
	if res.StatusCode == 200 {
		log.Println(targetUrl.String())
		robots, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", robots)
	} else {
		log.Println("No exist path:", targetUrl.String())
	}
	done <- true
}

func TestHttp(t *testing.T) {
	activeThreads := 0
	maxThreads := 10
	done := make(chan bool)

	baseurl := "http://www.baidu.com"
	f, err := os.Open("./wordlist.txt")
	if err != nil {
		log.Println("Error open file:", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		go checkfile(baseurl, scanner.Text(), done)
		activeThreads++
		if activeThreads >= maxThreads {
			<-done
			activeThreads -= 1

		}
	}

	for activeThreads > 0 {
		<-done
		activeThreads -= 1
	}

}
