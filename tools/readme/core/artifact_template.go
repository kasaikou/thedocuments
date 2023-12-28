package core

import (
	"errors"
	"fmt"
	"path/filepath"
)

type DirectoryArtifactTemplate struct {
	Title       string
	Path        PathTemplate
	Description string
	Content     []ArtifactTemplate
}

type ArtifactTemplate struct {
	Title       string
	FileType    string
	Path        PathTemplate
	PathKV      map[string]PathTemplate
	Info        string
	Description string
}

type PathTemplate struct {
	Basename string
	RelPath  string
	AbsPath  string
}

func FromAbsPath(root, abspath string) (PathTemplate, error) {
	path := PathTemplate{
		AbsPath:  abspath,
		Basename: filepath.Base(abspath),
	}
	if rel, err := filepath.Rel(root, abspath); err != nil {
		return path, errors.Join(fmt.Errorf("create relative filepath (base: %s, abs: %s)", root, abspath), err)
	} else {
		path.RelPath = rel
	}

	return path, nil
}
