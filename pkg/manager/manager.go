package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Model struct {
	Name     string   `json:"name"`
	Tag      string   `json:"tag"`
	Base     string   `json:"base"`
	Manifest Manifest `json:"manifest"`
}

type Manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int64  `json:"size"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int64  `json:"size"`
	} `json:"layers"`
}

type Blob struct {
	Digest string
	Size   int64
}

// NewModel and parses the model name and handles the follwing formats:
// <model> (e.g., "deepseek-r1")
// <model:tag> (e.g., "deepseek-r1:14b")
// In the first case "latest" will be the implied tag
func NewModel(model, base string) (*Model, error) {
	before, after, found := strings.Cut(model, ":")

	if found && after == "" {
		return nil, errors.New("tag syntax used without tag")
	}

	if before == "" {
		return nil, errors.New("invalid model name")
	}

	if !found {
		after = "latest"
	}

	return &Model{
		Name: before,
		Tag:  after,
		Base: base,
	}, nil
}

func (m Model) Pull() error {
	manifestFile := m.manifestFile()
	if _, err := os.Stat(manifestFile); !os.IsNotExist(err) {
		slog.Warn("manifest file already exists", "file", manifestFile)
		return nil
	}

	url := m.manifestURL()
	slog.Info("downloading model manifest", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download model manifest: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&m.Manifest); err != nil {
		return fmt.Errorf("failed to decode manifest: %w", err)
	}

	if err := m.downloadBlobs(); err != nil {
		return fmt.Errorf("failed to download blobs: %w", err)
	}

	file, err := os.Create(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(m.Manifest)
}

func (m Model) manifestFile() string {
	return filepath.Join(m.Base, "manifests", "registry.ollama.ai", "library", m.Name, m.Tag)
}

func (m Model) manifestURL() string {
	return fmt.Sprintf("https://registry.ollama.ai/v2/library/%s/manifests/%s", m.Name, m.Tag)
}

func (m Model) downloadBlobs() error {
	// blobs = layers + config
	blobs := make([]Blob, 0, len(m.Manifest.Layers)+1)
	for _, layer := range m.Manifest.Layers {
		blobs = append(blobs, Blob{
			Digest: layer.Digest,
			Size:   layer.Size,
		})
	}
	blobs = append(blobs, Blob{
		Digest: m.Manifest.Config.Digest,
		Size:   m.Manifest.Config.Size,
	})

	group := new(errgroup.Group)
	group.SetLimit(len(blobs))

	for _, blob := range blobs {
		group.Go(func() error {
			slog.Info("downloading blob", "blob", blob.Digest)
			url := fmt.Sprintf("https://registry.ollama.ai/v2/library/%s/blobs/%s", m.Name, blob.Digest)
			resp, err := http.Get(url)
			if err != nil {
				return nil
			}
			defer resp.Body.Close()

			path := filepath.Join(m.Base, "blobs", blob.Digest)
			file, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to download blob: %w", err)
			}
			defer file.Close()

			if err := file.Truncate(blob.Size); err != nil {
				slog.Warn("failed to truncate file", "blob", blob.Digest)
			}

			slog.Info("writing blob to disk", "blob", blob.Digest)
			if _, err := io.Copy(file, resp.Body); err != nil {
				return fmt.Errorf("failed to write blob to file: %w", err)
			}
			return nil
		})
	}
	return group.Wait()
}
