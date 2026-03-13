package config

import (
	"os"
	"sync"
)

// RuntimeConfig holds the runtime configuration set via UI or env vars.
type RuntimeConfig struct {
	BaseURL  string
	Token    string
	UserName string
}

var (
	mu      sync.RWMutex
	current = &RuntimeConfig{
		BaseURL:  os.Getenv("BACKEND_BASE_URL"),
		Token:    os.Getenv("BACKEND_JWT_TOKEN"),
		UserName: os.Getenv("BACKEND_USERNAME"),
	}
)

func Get() RuntimeConfig {
	mu.RLock()
	defer mu.RUnlock()
	return *current
}

func Set(cfg RuntimeConfig) {
	mu.Lock()
	defer mu.Unlock()
	current = &cfg
}
