package main

import (
	"log"
	"log/slog"
	"os"
	"runtime/pprof"
)

//go:generate go tool argsgen

type options struct {
	model     string `arg:"name of model to be downloaded,positional"`
	directory string `arg:"models diretory (OLLAMA_MODELS),required"`
	profile   bool   `arg:"create CPU profile"`
}

func main() {
	o := options{
		directory: os.Getenv("OLLAMA_MODELS"),
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
