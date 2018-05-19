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
		"https://www.ifa.plus",
		"https://b1.ifa.plus",
		"https://b2.ifa.plus",
		"https://b3.ifa.plus",
		"https://b4.ifa.plus",
		"https://b5.ifa.plus",
		"https://b6.ifa.plus",
		"https://b7.ifa.plus",
		"https://b8.ifa.plus",
		"https://b9.ifa.plus",
		"https://b10.ifa.plus",
		"https://b11.ifa.plus",
		"https://b12.ifa.plus"}
)

func get(url string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

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
			get(Host[0])
		}
	}()
	<-c
}
