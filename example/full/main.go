package main

import (
	"../../goquery"
	// "bytes"
	// "flag"
	"fmt"
	// "log"
	"net/http"
	// "net/url"
	// "runtime"
	// "strings"
	// "sync"
	"time"
)

var (
	Host = []string{
		"https://www.ifa.plus/pc/",
		"https://b1.ifa.plus/pc/",
		"https://b2.ifa.plus/pc/",
		"https://b3.ifa.plus/pc/",
		"https://b4.ifa.plus/pc/",
		"https://b5.ifa.plus/pc/",
		"https://b6.ifa.plus/pc/",
		"https://b7.ifa.plus/pc/",
		"https://b8.ifa.plus/pc/",
		"https://b9.ifa.plus/pc/",
		"https://b10.ifa.plus/pc/",
		"https://b11.ifa.plus/pc/",
		"https://b12.ifa.plus/pc/"}
)

func get(url string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %s %d %s", url, res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(doc)
	// Find the review items
	// doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	band := s.Find("a").Text()
	// 	title := s.Find("i").Text()
	// 	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// })
}
func main() {
	var ticker *time.Ticker = time.NewTicker(time.Duration(10) * time.Second)
	c := make(chan int, 1)
	go func() {
		for _ = range ticker.C {
			for _, u := range Host {
				get(u)
			}
		}
	}()
	<-c
}
