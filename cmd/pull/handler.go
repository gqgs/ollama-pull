package main

import (
	"fmt"

	"github.com/gqgs/ollama-pull/pkg/manager"
)

func handler(o options) error {
	model, err := manager.NewModel(o.model, o.directory)
	if err != nil {
		return fmt.Errorf("failed initializing model: %w", err)
	}

	if err := model.Pull(); err != nil {
		return fmt.Errorf("failed pulling the model: %w", err)
	}

	return nil
}
