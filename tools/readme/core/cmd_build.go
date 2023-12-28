package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"gopkg.in/yaml.v3"
)

type BuildConfig struct {
	configExt    string
	rootAbsDir   string
	configReader io.Reader
}

func Build(ctx context.Context, config BuildConfig) error {

	var decoder interface{ Decode(any) error }
	switch config.configExt {
	case "yml", "yaml":
		decoder = yaml.NewDecoder(config.configReader)
	case "json":
		decoder = json.NewDecoder(config.configReader)
	default:
		return fmt.Errorf("unknown file ext '%s'", config.configExt)
	}

	var configBody ArtifactConfig
	if err := decoder.Decode(&configBody); err != nil {
		return errors.Join(fmt.Errorf("unmarshal with format '%s'"), err)
	}

	for i := range configBody.Objects {
		configBody.Objects[i].Sort("path")
	}

	configBody.Sort()

	mapDirConfig := map[string]DirectoryArtifact{}
	for _, object := range configBody.Objects {
		filename := object.Path
		if filename == "" {
			filename = object.Files[0].Path
		}
		filename = filepath.Join(config.rootAbsDir, filename)
		stat, err := os.Stat(filename)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot found file '%s'", filename), err)
		}

		object.Path = filename

		if stat.IsDir() {
			dirname := filename
			artifact, exist := mapDirConfig[dirname]
			if !exist {
				artifact = DirectoryArtifact{Path: dirname}
			}
			artifact.DirectoryContent = object
			mapDirConfig[dirname] = artifact

		} else {
			dirname := filepath.Dir(filename)
			artifact, exist := mapDirConfig[dirname]
			if !exist {
				artifact = DirectoryArtifact{Path: dirname}
			}
			artifact.Content = append(artifact.Content, object)
			mapDirConfig[dirname] = artifact
		}

	}

	wg := sync.WaitGroup{}
	template.New("root").Parse(filename)
	for _, artifact := range mapDirConfig {
		go func(artifact DirectoryArtifact) {
			defer wg.Done()

		}(artifact)
	}

}
