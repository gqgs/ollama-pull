package downloader

import "errors"

type Downloader interface {
	Download(url, output string) error
}

func New(name string) (Downloader, error) {
	switch name {
	case "aria":
		return NewAria()
	case "http":
		return NewHttp(), nil
	default:
		return nil, errors.New("invlid download option:" + name)
	}
}
