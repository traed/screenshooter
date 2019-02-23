package main

import (
	"flag"
	"log"
	"net/url"
	"sync"

	"../../pkg/capturer"
	"../../pkg/chrome"
)

func main() {
	path := flag.String("p", "../uploads", "The place where screenshots are saved.")
	flag.Parse()

	var wg sync.WaitGroup
	browser := chrome.NewChrome()
	browser.SetScreenshotPath(*path)

	for _, arg := range flag.Args() {
		u, err := url.ParseRequestURI(arg)
		if err != nil {
			log.Print("Invalid URL specified")
			return
		}

		wg.Add(1)
		capturer := capturer.Capturer{Browser: browser, URL: u, WG: &wg}
		go capturer.Execute()
	}

	wg.Wait()
}
