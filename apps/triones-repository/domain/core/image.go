package core

import (
	"time"

	"github.com/hashicorp/go-version"
)

type Image struct {
	Version     version.Version `json:"version"`
	Tags        []string        `json:"tags"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	FilePath    string          `json:"filepath"`
	CreateTime  time.Time       `json:"create_time"`
}
