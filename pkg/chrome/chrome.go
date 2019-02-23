package chrome

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Chrome contains information about a Google Chrome
// instance, with methods to run on it.
type Chrome struct {
	Resolution    string
	ChromeTimeout int
	Path          string
	UserAgent     string

	ScreenshotPath string
}

// NewChrome - Factory function
func NewChrome() *Chrome {
	chrome := new(Chrome)
	chrome.Path = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	chrome.Resolution = "1024x768"
	chrome.ChromeTimeout = 10
	chrome.UserAgent = "Screenshooter"

	return chrome
}

// SetScreenshotPath sets the path for screenshots
func (chrome *Chrome) SetScreenshotPath(p string) error {
	p, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("Destination path does not exist")
	}

	// Log
	chrome.ScreenshotPath = p

	return nil
}

// ScreenshotURL takes a screenshot of a URL
func (chrome *Chrome) ScreenshotURL(targetURL *url.URL, destination string) {
	// Basic arguments for headless chrome
	var chromeArguments = []string{
		"--headless", "--disable-gpu", "--hide-scrollbars",
		"--disable-crash-reporter",
		"--user-agent=" + chrome.UserAgent,
		"--window-size=" + chrome.Resolution, "--screenshot=" + destination,
	}

	// Handle 'cant run as root'.
	if os.Geteuid() == 0 {
		chromeArguments = append(chromeArguments, "--no-sandbox")
	}

	chromeArguments = append(chromeArguments, targetURL.String())

	if targetURL.Scheme != "https" {
		log.Print("Skipping non-encrypted URLs.")
		return
	}

	// Get a context to run the command in
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(chrome.ChromeTimeout)*time.Second)
	defer cancel()

	// Prepare the command to run...
	cmd := exec.CommandContext(ctx, chrome.Path, chromeArguments...)

	// ... and run it!
	if err := cmd.Start(); err != nil {
		log.Printf("An error occurred while starting chrome: %s", err)
	}

	// Wait for the screenshot to finish and handle the error that may occur.
	if err := cmd.Wait(); err != nil {

		// If if this error was as a result of a timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Screenshot timed out: %s", err)
			return
		}

		log.Printf("An error occurred while taking a screenshot: %s", err)
		return
	}

	log.Printf("Screenshot taken from %s", targetURL.String())
}
