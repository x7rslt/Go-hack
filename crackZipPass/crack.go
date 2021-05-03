package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/alexmullins/zip"
)

var (
	zipfile, wordlist string
	speed             int
)

func init() {
	flag.StringVar(&zipfile, "z", "", "input zipfile")
	flag.StringVar(&wordlist, "w", "", "input password list")
	flag.IntVar(&speed, "s", runtime.NumCPU(), "default use 4 cpu work")
}

func main() {
	flag.Parse()
	fmt.Println(zipfile, wordlist, speed)

	//判断输入zip文件、字典、速度参数
	if zipfile == "" || wordlist == "" {
		fmt.Println("Please input your zipfile or wordlist.")
		os.Exit(0)
	}

	word := make(chan string, 0) //word 没有接收，为什么没有进程堵塞？
	found := make(chan string, 0)

	var wg sync.WaitGroup

	for i := 0; i <= speed; i++ {
		wg.Add(1)
		go Crackzip(word, found, &wg)
	}

	//读取wordlist
	wordfile, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer wordfile.Close()

	scanner := bufio.NewScanner(wordfile)
	go func() {
		for scanner.Scan() {
			pass := scanner.Text()
			word <- pass
		}
		close(word)
	}()

	done := make(chan bool)

	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case f := <-found:
		println("[+]Found Password")
		println("[+]Password is :", f)
		return
	case <-done:
		println("[+] Password not found")
		return
	}
}

func Crackzip(word <-chan string, found chan<- string, wg *sync.WaitGroup) {
	zipr, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer zipr.Close()
	defer wg.Done()
	fmt.Println("word :", word)
	for w := range word {
		fmt.Println(w)
		for _, z := range zipr.File {
			z.SetPassword(w)
			_, err := z.Open()
			if err == nil {
				found <- w
			}
		}
	}

}
