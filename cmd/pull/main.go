package main

import (
	"log/slog"
	"os"
)

//go:generate go run github.com/gqgs/argsgen@latest

type options struct {
	model     string `arg:"name of model to be downloaded,positional"`
	directory string `arg:"models diretory,required"`
}

func main() {
	o := options{
		directory: os.Getenv("OLLAMA_MODELS"),
	}
	o.MustParse()

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
