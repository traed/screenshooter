package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"../../pkg/capturer"
	"../../pkg/chrome"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage:")
		fmt.Fprintln(flag.CommandLine.Output(), "scsh [-d=<path>] <url> [<url>...]")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	path := flag.String("d", "../uploads", "Directory where screenshots are saved.")
	flag.Parse()

	var wg sync.WaitGroup
	browser := chrome.NewChrome()
	browser.SetScreenshotPath(*path)

	urls := flag.Args()

	if len(urls) == 0 {
		log.Fatal("No URLs provided.")
	}

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
