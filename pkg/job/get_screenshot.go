package job

import (
	"image/png"
	"io"
	"os"
)

// GetScreenshotJob is an impementation of worker.Job interface
type GetScreenshotJob struct {
	pw       *io.PipeWriter
	filepath string
}

// NewGetScreenshotJob constructs a getScreenshotJob
func NewGetScreenshotJob(pw *io.PipeWriter, fp string) GetScreenshotJob {
	return GetScreenshotJob{pw: pw, filepath: fp}
}

// Execute a GetScreenshotJob
func (j *GetScreenshotJob) Execute() error {
	defer j.pw.Close()

	_, err := os.Stat(j.filepath)
	if os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}

	file, err := os.Open(j.filepath)
	defer file.Close()
	if err != nil {
		return err
	}

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	err = png.Encode(j.pw, img)
	if err != nil {
		return err
	}

	return nil
}
