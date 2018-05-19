package main

import (
	"./colly"
	"fmt"
	"net/http"
	"time"
)

var (
	// Host = []string{
	// 	"https://www.ifa.plus/pc/",
	// 	"https://b1.ifa.plus/pc/",
	// 	"https://b2.ifa.plus/pc/",
	// 	"https://b3.ifa.plus/pc/",
	// 	"https://b4.ifa.plus/pc/",
	// 	"https://b5.ifa.plus/pc/",
	// 	"https://b6.ifa.plus/pc/",
	// 	"https://b7.ifa.plus/pc/",
	// 	"https://b8.ifa.plus/pc/",
	// 	"https://b9.ifa.plus/pc/",
	// 	"https://b10.ifa.plus/pc/",
	// 	"https://b11.ifa.plus/pc/",
	// 	"https://b12.ifa.plus/pc/"}
	Host = []string{
		// "https://en.wikipedia.org",
		"http://www.91lym.com/"}
)

func start(url string) {
	// Instantiate default collector
	c := colly.NewCollector(
	// MaxDepth is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	// colly.MaxDepth(10),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	c.MaxDepth = 10
	c.IgnoreRobotsTxt = true
	// c.SetCookies()
	// sto := c.GetStorage()
	// sto.SetCookies(url, "connect.sid=s%3AwnC_5kcvATUaATV_qeEC-A0Ofb4-P9PY.TZmZCMRd%2BvoVzdKCQ3eO8tk7%2FMt1GN35iGx2JaC22tU")
	// c.SetStorage(sto)
	COOKIE_MAX_MAX_AGE := time.Hour * 24 / time.Second // 单位：秒。
	maxAge := int(COOKIE_MAX_MAX_AGE)
	uid := "s%3AwnC_5kcvATUaATV_qeEC-A0Ofb4-P9PY.TZmZCMRd%2BvoVzdKCQ3eO8tk7%2FMt1GN35iGx2JaC22tU"

	uid_cookie := &http.Cookie{
		Name:     "connect.sid",
		Value:    uid,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   maxAge,
	}
	c.SetCookies(url, []*http.Cookie{uid_cookie})
	// On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	fmt.Println(link)
	// 	// Visit link found on page
	// 	e.Request.Visit(link)
	// })
	c.Limit(&colly.LimitRule{
		// Parallelism: 2,
		RandomDelay: 10 * time.Second,
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println("e:", e)
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
		// link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		e.ForEach("table tbody tr", func(_ int, el *colly.HTMLElement) {
			// mail := Mail{
			// 	Title:   el.ChildText("td:nth-of-type(1)"),
			// 	Link:    el.ChildAttr("td:nth-of-type(1)", "href"),
			// 	Author:  el.ChildText("td:nth-of-type(2)"),
			// 	Date:    el.ChildText("td:nth-of-type(3)"),
			// 	Message: el.ChildText("td:nth-of-type(4)"),
			// }
			// threads[threadSubject] = append(threads[threadSubject], mail)
			fmt.Printf("时间\t比赛名称\t4:4额度\n")
			fmt.Printf("%s\t%s\t%s\n", el.ChildText("td:time"), el.ChildText("td:time"), el.ChildText("td:four good"))
		})

		// c.Visit(e.Request.AbsoluteURL(link))

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
