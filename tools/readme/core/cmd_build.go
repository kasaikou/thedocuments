package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/kasaikou/thedocuments/tools"
	"gopkg.in/yaml.v3"
)

func newErrFailedBuild(configFilename string, includes ...error) error {
	err := fmt.Errorf("failed build from '%s'", configFilename)
	if len(includes) > 0 {
		err = errors.Join(append([]error{err}, includes...)...)
	}
	return err
}

func Build(ctx context.Context, configFilename string) error {

	logger := tools.LoggerFromContext(ctx)
	logger = logger.With(slog.String("configFile", configFilename))

	if !filepath.IsAbs(configFilename) {
		panic("configFilename argument must be absolute filepath")
	}

	configDir := filepath.Dir(configFilename)
	configExt := filepath.Ext(configFilename)
	configReader, err := os.Open(configFilename)

	if err != nil {
		logger.Error("Cannot open config file", slog.Any("error", err))
		return newErrFailedBuild(configFilename, err)
	}

	var decoder interface{ Decode(any) error }
	switch configExt {
	case ".yml", ".yaml":
		decoder = yaml.NewDecoder(configReader)
	case ".json":
		decoder = json.NewDecoder(configReader)
	default:
		logger.Error("Unknown file ext", slog.String("ext", configExt))
		return newErrFailedBuild(configFilename)
	}

	var config ArtifactConfig
	if err := decoder.Decode(&config); err != nil {
		logger.Error("Failed to unmarshal with format", slog.String("ext", configExt), slog.Any("error", err))
		return newErrFailedBuild(configFilename, err)
	}

	for i := range config.Objects {
		config.Objects[i].Sort("path")
	}

	config.Sort()

	mapDirConfig := map[string]DirectoryArtifactTemplate{}
	for _, object := range config.Objects {
		filename := object.Path
		if filename == "" {
			filename = object.Files[0].Path
		}
		filename = filepath.Join(configDir, filename)
		stat, err := os.Stat(filename)
		if err != nil {
			logger.Error("Could not find file", slog.String("path", filename), slog.Any("error", err))
			return newErrFailedBuild(configFilename, err)
		}

		object.Path = filename

		if stat.IsDir() {
			dirname := filename
			artifact, exist := mapDirConfig[dirname]
			if !exist {
				artifact = DirectoryArtifactTemplate{}
				if path, err := FromAbsPath(configDir, dirname); err != nil {
					logger.Error("Failed to create path object", slog.Any("error", err))
					return newErrFailedBuild(configFilename, err)
				} else {
					artifact.Path = path
				}
			}

			artifact.Description = object.Description
			mapDirConfig[dirname] = artifact

		} else {
			dirname := filepath.Dir(filename)
			artifact, exist := mapDirConfig[dirname]
			if !exist {
				artifact = DirectoryArtifactTemplate{}
				if path, err := FromAbsPath(configDir, dirname); err != nil {
					logger.Error("Failed to create path object", slog.Any("error", err))
					return newErrFailedBuild(configFilename, err)
				} else {
					artifact.Path = path
				}
			}

			objectArtifact := ArtifactTemplate{
				Title:       object.Title,
				Info:        object.Info,
				Description: object.Description,
			}

			if object.Path == "" {
				objectArtifact.FileType = "kv"
				for _, file := range object.Files {
					if path, err := FromAbsPath(configDir, filename); err != nil {
						logger.Error("Failed to create path object", slog.Any("error", err))
						return newErrFailedBuild(configFilename, err)
					} else {
						objectArtifact.PathKV[file.Key] = path
					}
				}
			} else {
				if path, err := FromAbsPath(configDir, filename); err != nil {
					logger.Error("Failed to create path object", slog.Any("error", err))
					return newErrFailedBuild(configFilename, err)
				} else {
					objectArtifact.FileType = "single"
					objectArtifact.Path = path
				}
			}

			artifact.Content = append(artifact.Content, objectArtifact)
			mapDirConfig[dirname] = artifact
		}

	}

	wg := sync.WaitGroup{}
	tmpl, err := template.New("root").Parse(config.Layout)
	if err != nil {
		logger.Error("Cannot load layout template", slog.Any("error", err))
		return newErrFailedBuild(configFilename, err)
	}

	wholeErr := []error{}
	lockWholeErr := sync.Mutex{}
	appendWholeError := func(err error) {
		lockWholeErr.Lock()
		defer lockWholeErr.Unlock()
		wholeErr = append(wholeErr, err)
	}

	logger.Info(fmt.Sprintf("Scheduled %d save artifacts", len(mapDirConfig)))

	for _, artifact := range mapDirConfig {
		wg.Add(1)
		go func(artifact DirectoryArtifactTemplate) {
			defer wg.Done()
			tmpl, err := tmpl.Clone()
			if err != nil {
				appendWholeError(err)
				return
			}

			file, err := os.Create(filepath.Join(artifact.Path.AbsPath, config.OutputFilename))
			if err != nil {
				appendWholeError(err)
				return
			}

			if err := tmpl.Execute(file, artifact); err != nil {
				appendWholeError(err)
				return
			}

		}(artifact)
	}
	wg.Wait()

	if len(wholeErr) > 0 {
		logger.Error("Failed in save artifacts", slog.Any("errors", wholeErr))
		return newErrFailedBuild(configFilename, wholeErr...)
	} else {
		return nil
	}
}
