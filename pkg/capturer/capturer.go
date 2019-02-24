package capturer

import (
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"../chrome"

	"github.com/parnurzeal/gorequest"
)

//Capturer - Struct that captures screenshots
type Capturer struct {
	URL     *url.URL
	Browser *chrome.Chrome
	Dest    string
}

// NewCapturer - Construct a new Capturer
func NewCapturer(imagePath string, url *url.URL) Capturer {
	browser := chrome.NewChrome()
	browser.SetScreenshotPath(imagePath)

	fname := SafeFileName(url.String())
	dest := filepath.Join(browser.ScreenshotPath, fname)

	capturer := Capturer{Browser: browser, URL: url, Dest: dest}

	return capturer
}

// Execute - Execute the screenshot capture
func (c *Capturer) Execute() {
	request := gorequest.
		New().
		Timeout(time.Duration(c.Browser.ChromeTimeout)*time.Second).
		Set("User-Agent", c.Browser.UserAgent)

	resp, _, errs := request.Get(c.URL.String()).End()
	if errs != nil {
		for _, err := range errs {
			log.Print(err.Error())
		}
		return
	}

	finalURL := resp.Request.URL

	c.Browser.ScreenshotURL(finalURL, c.Dest)
}

// GetFilename returns the formated filename used by the Capturer
func (c *Capturer) GetFilename() string {
	return filepath.Base(c.Dest)
}

// SafeFileName converts str into a filename with the
func SafeFileName(str string) string {
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

	return name + ".png"
}
