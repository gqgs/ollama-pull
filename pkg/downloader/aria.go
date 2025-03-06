package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type aria struct {
	aria2c string
}

func (a *aria) Download(url, output string) error {
	cmd := exec.Command(a.aria2c, url, "--async-dns=false", "--dir", filepath.Dir(output), "-o", filepath.Base(output))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewAria() (*aria, error) {
	aria2c, err := exec.LookPath("aria2c")
	if err != nil {
		return nil, fmt.Errorf("failed to find aria2c executable")
	}
	return &aria{
		aria2c,
	}, nil
}
