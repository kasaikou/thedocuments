package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"gopkg.in/yaml.v3"
)

func Build(ctx context.Context, configFilename string) error {

	if !filepath.IsAbs(configFilename) {
		panic("configFilename argument must be absolute filepath")
	}

	configDir := filepath.Dir(configFilename)
	configExt := filepath.Ext(configFilename)
	configReader, err := os.Open(configFilename)

	if err != nil {
		return errors.Join(fmt.Errorf("cannot open file '%s'", configFilename), err)
	}

	var decoder interface{ Decode(any) error }
	switch configExt {
	case ".yml", ".yaml":
		decoder = yaml.NewDecoder(configReader)
	case ".json":
		decoder = json.NewDecoder(configReader)
	default:
		return fmt.Errorf("unknown file ext '%s'", configExt)
	}

	var config ArtifactConfig
	if err := decoder.Decode(&config); err != nil {
		return errors.Join(fmt.Errorf("unmarshal with format '%s'", configExt), err)
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
			return errors.Join(fmt.Errorf("cannot found file '%s'", filename), err)
		}

		object.Path = filename

		if stat.IsDir() {
			dirname := filename
			artifact, exist := mapDirConfig[dirname]
			if !exist {
				artifact = DirectoryArtifactTemplate{}
				if path, err := FromAbsPath(configDir, dirname); err != nil {
					return err
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
					return err
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
						return err
					} else {
						objectArtifact.PathKV[file.Key] = path
					}
				}
			} else {
				if path, err := FromAbsPath(configDir, filename); err != nil {
					return err
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
		return errors.Join(errors.New("cannot load layout template"), err)
	}

	wholeErr := []error{}
	lockWholeErr := sync.Mutex{}
	appendWholeError := func(err error) {
		lockWholeErr.Lock()
		defer lockWholeErr.Unlock()
		wholeErr = append(wholeErr, err)
	}

	log.Println("Scheduled", len(mapDirConfig), "save artifacts")

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
		return errors.Join(wholeErr...)
	} else {
		return nil
	}
}
