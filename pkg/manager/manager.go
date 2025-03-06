package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gqgs/ollama-pull/pkg/downloader"
	"golang.org/x/sync/errgroup"
)

const (
	registry = "registry.ollama.ai"
)

type Model struct {
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Base     string `json:"base"`
	Manifest struct {
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
	} `json:"manifest"`
	downloader downloader.Downloader
}

type blob struct {
	Digest string
	Size   int64
}

// NewModel and parses the model name and handles the following formats:
// <model> (e.g., "deepseek-r1")
// <model:tag> (e.g., "deepseek-r1:14b")
// In the first case "latest" will be the implied tag
func NewModel(model, base string, downloader downloader.Downloader) (*Model, error) {
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
		Name:       before,
		Tag:        after,
		Base:       base,
		downloader: downloader,
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

	if err := os.MkdirAll(filepath.Dir(manifestFile), os.ModePerm); err != nil {
		return fmt.Errorf("failed creating directory: %w", err)
	}

	file, err := os.Create(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(m.Manifest)
}

func (m Model) manifestFile() string {
	return filepath.Join(m.Base, "manifests", registry, "library", m.Name, m.Tag)
}

func (m Model) manifestURL() string {
	return fmt.Sprintf("https://%s/v2/library/%s/manifests/%s", registry, m.Name, m.Tag)
}

func (m Model) downloadBlobs() error {
	// blobs = layers + config
	blobs := make([]blob, 0, len(m.Manifest.Layers)+1)
	for _, layer := range m.Manifest.Layers {
		blobs = append(blobs, blob{
			Digest: layer.Digest,
			Size:   layer.Size,
		})
	}
	blobs = append(blobs, blob{
		Digest: m.Manifest.Config.Digest,
		Size:   m.Manifest.Config.Size,
	})

	group := new(errgroup.Group)
	group.SetLimit(len(blobs))

	for _, blob := range blobs {
		group.Go(func() error {
			slog.Info("downloading blob", "blob", blob.Digest)

			url := fmt.Sprintf("https://%s/v2/library/%s/blobs/%s", registry, m.Name, blob.Digest)

			path := filepath.Join(m.Base, "blobs", blob.Digest)
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return fmt.Errorf("failed creating directory: %w", err)
			}

			return m.downloader.Download(url, path)
		})
	}
	return group.Wait()
}
