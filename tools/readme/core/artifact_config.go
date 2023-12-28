package core

import (
	"fmt"
	"slices"
	"strings"
)

type ArtifactConfig struct {
	Layout  string           `yaml:"layout"`
	Objects []ArtifactObject `yaml:"objects"`
}

type ArtifactObject struct {
	Title       string         `yaml:"title"`
	Info        string         `yaml:"info"`
	Path        string         `yaml:"path"`
	Files       []ArtifactFile `yaml:"files"`
	Description string         `yaml:"description"`
}

type ArtifactFile struct {
	Key  string `yaml:"key"`
	Path string `yaml:"path"`
}

func (a *ArtifactObject) Sort(by string) {
	switch by {
	case "key":
		if len(a.Files) > 1 {
			slices.SortFunc(a.Files, func(a, b ArtifactFile) int {
				return strings.Compare(a.Key, b.Key)
			})
		}
	case "path":
		if len(a.Files) > 1 {
			slices.SortFunc(a.Files, func(a, b ArtifactFile) int {
				return strings.Compare(a.Path, b.Path)
			})
		}
	default:
		panic(fmt.Sprintf("unknown argument: sort by '%s'", by))
	}
}

func (a *ArtifactConfig) Sort() {
	slices.SortFunc(a.Objects, func(l, r ArtifactObject) int {
		left := l.Path
		right := r.Path
		if left == "" {
			left = l.Files[0].Path
		}
		if right == "" {
			right = r.Files[0].Path
		}
		return strings.Compare(left, right)
	})
}
