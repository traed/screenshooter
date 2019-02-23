package main

import (
	"log"
	"net/url"
	"os"
	"sync"

	"../../pkg/capturer"
	"../../pkg/chrome"
)

func main() {
	args := os.Args[1:]
	var wg sync.WaitGroup
	browser := chrome.NewChrome()

	for _, arg := range args {
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
