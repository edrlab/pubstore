package view

import (
	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

type View struct {
	*conf.Config
	*stor.Store
}

func Init(c *conf.Config, s *stor.Store) View {
	return View{
		Config: c,
		Store:  s,
	}
}

// contentTypeToFormat utility
func contentTypeToFormat(contentType string) string {
	switch contentType {
	case "application/epub+zip":
		return "epub"
	case "application/pdf+lcp":
		return "pdf"
	case "application/audiobook+lcp":
		return "audiobook"
	case "application/divina+lcp":
		return "divina"
	default:
		return "unknown"
	}
}

// formatToContentType utility
func formatToContentType(format string) string {
	switch format {
	case "epub":
		return "application/epub+zip"
	case "pdf":
		return "application/pdf+lcp"
	case "audiobook":
		return "application/audiobook+lcp"
	case "divina":
		return "application/divina+lcp"
	default:
		return "unknown"
	}
}
