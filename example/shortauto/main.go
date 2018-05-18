package main

import (
	"fmt"
	"net/http"
	"time"

	"../../fetchbot"
)

func main() {
	f := fetchbot.New(fetchbot.HandlerFunc(handler))
	f.AutoClose = true
	f.WorkerIdleTTL = time.Second
	queue := f.Start()
	queue.SendStringHead("https://www.ifa.plus", "https://b1.ifa.plus/")
	queue.Block()
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
}
