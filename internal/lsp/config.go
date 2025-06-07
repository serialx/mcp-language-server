package lsp

import (
	"encoding/json"
	"fmt"
	"os"
)

// ServerConfig holds language server configuration
type ServerConfig struct {
	// InitializationOptions are passed to the language server during initialization
	InitializationOptions map[string]any `json:"initializationOptions,omitempty"`

	// Settings are returned when the server requests workspace/configuration
	Settings map[string]any `json:"settings,omitempty"`
}

// LoadServerConfig loads configuration from a JSON file
func LoadServerConfig(path string) (*ServerConfig, error) {
	if path == "" {
		// Return empty config if no path provided
		return &ServerConfig{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetDefaultConfig returns default configurations for known language servers
func GetDefaultConfig(lspCommand string) *ServerConfig {
	switch lspCommand {
	case "gopls":
		return &ServerConfig{
			InitializationOptions: map[string]any{
				"codelenses": map[string]bool{
					"generate":           true,
					"regenerate_cgo":     true,
					"test":               true,
					"tidy":               true,
					"upgrade_dependency": true,
					"vendor":             true,
					"vulncheck":          false,
				},
			},
		}
	default:
		// Return empty config for other language servers
		return &ServerConfig{}
	}
}
