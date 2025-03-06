package main

import (
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime/pprof"
)

//go:generate go tool argsgen

type options struct {
	model      string `arg:"name of model to be downloaded,positional"`
	directory  string `arg:"models diretory (OLLAMA_MODELS),required"`
	downloader string `arg:"downloader type (aria|http)"`
	profile    bool   `arg:"create CPU profile"`
}

func main() {
	o := options{
		directory:  os.Getenv("OLLAMA_MODELS"),
		downloader: "http",
	}

	// default to aria if it exists
	if _, err := exec.LookPath("aria2c"); err == nil {
		o.downloader = "aria"
	}

	o.MustParse()

	if o.profile {
		f, err := os.Create("default.pgo")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("error starting CPU profile", "err", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
