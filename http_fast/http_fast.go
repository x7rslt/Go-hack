package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Visit(url string, done chan bool) {
	client := http.Client{Timeout: time.Duration(8 * time.Second)}
	res, err := client.Get(url)
	if err != nil {
		log.Println("Error Client Get :", err)
	} else {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println("Goquery Error:", err)
		} else {
			title := doc.Find("title").Text()
			log.Println("Visit domain:", title)
		}

	}

	done <- true
}

func main() {
	activeThreads := 0
	maxThreads := 100000
	done := make(chan bool)
	f, err := os.Open("./domains.txt")
	if err != nil {
		log.Println("File open error :", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		go Visit("http://"+scanner.Text(), done)
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
	fmt.Println("All domain visited.")
}
