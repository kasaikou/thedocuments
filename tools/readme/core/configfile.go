package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kasaikou/thedocuments/tools"
	"gopkg.in/yaml.v3"
)

func LoadAndFormat(ctx context.Context, configFilename string) (ArtifactConfig, error) {
	logger := tools.LoggerFromContext(ctx)
	logger = logger.With(slog.String("configFile", configFilename))

	if !filepath.IsAbs(configFilename) {
		panic("configFilename argument must be absolute filepath")
	}

	configExt := filepath.Ext(configFilename)
	configReader, err := os.Open(configFilename)

	if err != nil {
		return ArtifactConfig{}, errors.Join(fmt.Errorf("cannot open config file '%s'", configFilename), err)
	}

	var decoder interface{ Decode(any) error }
	switch configExt {
	case ".yml", ".yaml":
		decoder = yaml.NewDecoder(configReader)
	case ".json":
		decoder = json.NewDecoder(configReader)
	default:
		return ArtifactConfig{}, fmt.Errorf("unknown file ext: '%s'", configExt)
	}

	var config ArtifactConfig
	if err := decoder.Decode(&config); err != nil {
		return ArtifactConfig{}, errors.Join(fmt.Errorf("failed to unmarshal with format '%s'", configExt), err)
	}

	for i := range config.Objects {
		config.Objects[i].Sort("path")
	}

	config.Sort()

	return config, nil
}
