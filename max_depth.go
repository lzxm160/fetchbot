package main

import (
	"./colly"
	"fmt"
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

func start(url string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(10),
	)

	// On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	fmt.Println(link)
	// 	// Visit link found on page
	// 	e.Request.Visit(link)
	// })
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 10 * time.Second,
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println("e:", e)
		// t := make([]transcript, 0)
		// e.ForEach(".topic-media-row", func(_ int, el *colly.HTMLElement) {
		// 	t = append(t, transcript{
		// 		Speaker: el.ChildText(".speaker-label"),
		// 		Text:    el.ChildText(".transcript-text-block"),
		// 	})
		// })
		// jsonData, err := json.MarshalIndent(t, "", "  ")
		// if err != nil {
		// 	return
		// }
		// ioutil.WriteFile(colly.SanitizeFileName(e.Request.Ctx.Get("date")+"_"+e.Request.Ctx.Get("slug"))+".json", jsonData, 0644)
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		c.Visit(e.Request.AbsoluteURL(link))

	})
	//lianjie body
	// Start scraping on https://en.wikipedia.org
	c.Visit(url)
	// c.Visit("https://en.wikipedia.org")
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())

	})
	c.Wait()
}
func main() {
	var ticker *time.Ticker = time.NewTicker(time.Duration(10) * time.Second)
	c := make(chan int, 1)
	go func() {
		for _ = range ticker.C {
			for _, u := range Host {
				start(u)
			}
		}
	}()
	<-c
}
