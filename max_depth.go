package main

import (
	"fmt"

	"./colly"
)

func main() {
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
		fmt.Println("e:", e)
	})
	//lianjie body
	// Start scraping on https://en.wikipedia.org
	c.Visit("https://b12.ifa.plus/pc/")
	// c.Visit("https://en.wikipedia.org")
}
