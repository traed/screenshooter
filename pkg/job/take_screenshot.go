package job

import (
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/traed/screenshooter/pkg/chrome"
	"github.com/traed/screenshooter/pkg/util"
)

//TakeScreenshotJob - Struct that captures screenshots
type TakeScreenshotJob struct {
	URL     *url.URL
	Browser *chrome.Chrome
	Dest    string
}

// NewTakeScreenshotJob - Construct a new TakeScreenshotJob
func NewTakeScreenshotJob(imagePath string, url *url.URL) TakeScreenshotJob {
	browser := chrome.NewChrome()
	browser.SetScreenshotPath(imagePath)

	fname := util.SafeFileName(url.String())
	dest := filepath.Join(browser.ScreenshotPath, fname)

	return TakeScreenshotJob{Browser: browser, URL: url, Dest: dest}
}

// Execute - Execute the screenshot capture
func (j *TakeScreenshotJob) Execute() error {
	client := &http.Client{
		Timeout: time.Duration(j.Browser.ChromeTimeout) * time.Second,
	}
	req, err := http.NewRequest("GET", j.URL.String(), nil)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	req.Header.Set("User-Agent", j.Browser.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	finalURL := resp.Request.URL

	j.Browser.ScreenshotURL(finalURL, j.Dest)

	return nil
}

// GetFilename returns the formated filename used by the TakeScreenshotJob
func (j *TakeScreenshotJob) GetFilename() string {
	return filepath.Base(j.Dest)
}
