package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func NearestFilename(absFromDir, basename string) (dir, filename string, err error) {

	if !filepath.IsAbs(absFromDir) {
		panic("absFromDir arguments must be absolute directory path")
	}

	absDir := absFromDir

	for {
		filename := filepath.Join(absDir, basename)
		if _, err := os.Stat(filename); err == nil {
			return absDir, filename, nil
		} else if parent := filepath.Dir(absDir); parent == absFromDir {
			return absDir, filename, fmt.Errorf("cannot found '%s' in '%s' and its parent directories", basename, absFromDir)
		} else {
			absDir = parent
		}
	}
}
