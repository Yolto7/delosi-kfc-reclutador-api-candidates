package config

import (
	"fmt"
	"os"

	"github.com/Yolto7/api-candidates/internal/domain/config"
)

func Load() (*config.Config, error) {
  cfg := &config.Config{}

  // --- Dynamo ---
  cfg.CANDIDATES_TABLE_NAME = os.Getenv("CANDIDATES_TABLE_NAME")
  if cfg.CANDIDATES_TABLE_NAME == "" {
    return nil, fmt.Errorf("CANDIDATES_TABLE_NAME environment variable is empty")
  }
   
  // --- Return ---
  return cfg, nil
}

