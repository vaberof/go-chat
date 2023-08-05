package logs

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

type Logs struct {
	root *zap.SugaredLogger
}

func New(configPath string) (*Logs, error) {
	logger, err := getLogger(configPath)
	if err != nil {
		return nil, err
	}
	return &Logs{root: logger}, nil
}

func (logs *Logs) WithName(name string) *zap.SugaredLogger {
	return logs.root.With().Named(name)
}

func getLogger(configPath string) (*zap.SugaredLogger, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open logger config")
	}

	configContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read logger config")
	}

	var cfg zap.Config
	if err = json.Unmarshal(configContent, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logger config content %s", err)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger %s", err)
	}

	return logger.Sugar(), nil
}
