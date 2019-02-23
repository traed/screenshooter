package capturer

import (
	"crypto/tls"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"../chrome"

	"github.com/parnurzeal/gorequest"
)

//Capturer - Struct that captures screenshots
type Capturer struct {
	URL     *url.URL
	Browser *chrome.Chrome
	WG      *sync.WaitGroup
}

// Execute - Execute the screenshot capture
func (c *Capturer) Execute() {
	defer c.WG.Done()

	request := gorequest.
		New().
		Timeout(time.Duration(c.Browser.ChromeTimeout)*time.Second).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("User-Agent", c.Browser.UserAgent)

	resp, _, errs := request.Get(c.URL.String()).End()
	if errs != nil {
		for _, err := range errs {
			log.Print(err.Error())
		}
		return
	}

	finalURL := resp.Request.URL
	fname := safeFileName(c.URL.String()) + ".png"
	dest := filepath.Join(c.Browser.ScreenshotPath, fname)

	c.Browser.ScreenshotURL(finalURL, dest)
	log.Printf("Saved screenshot to %s", dest)
}

func safeFileName(str string) string {
	name := strings.ToLower(str)
	name = strings.Trim(name, " ")

	separators, err := regexp.Compile(`[ &_=+:]`)
	if err == nil {
		name = separators.ReplaceAllString(name, "-")
	}

	legal, err := regexp.Compile(`[^[:alnum:]-.]`)
	if err == nil {
		name = legal.ReplaceAllString(name, "")
	}

	for strings.Contains(name, "--") {
		name = strings.Replace(name, "--", "-", -1)
	}

	return name
}
