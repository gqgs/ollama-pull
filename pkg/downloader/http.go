package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type httpDownloader struct{}

func (httpDownloader) Download(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download blob: %w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to download blob: %w", err)
	}
	defer file.Close()

	if err := file.Truncate(resp.ContentLength); err != nil {
		slog.Warn("failed to truncate file")
	}

	slog.Info("writing blob to disk", "url", url)
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write blob to file: %w", err)
	}
	return nil
}

func NewHttp() *httpDownloader {
	return new(httpDownloader)
}
