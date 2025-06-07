package lsp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServerConfig(t *testing.T) {
	tests := []struct {
		name        string
		configJSON  string
		expectError bool
		validate    func(t *testing.T, config *ServerConfig)
	}{
		{
			name:       "empty config file path",
			configJSON: "",
			validate: func(t *testing.T, config *ServerConfig) {
				if config == nil {
					t.Fatal("expected non-nil config")
				}
				if config.InitializationOptions != nil {
					t.Error("expected nil InitializationOptions")
				}
				if config.Settings != nil {
					t.Error("expected nil Settings")
				}
			},
		},
		{
			name: "valid config with initialization options",
			configJSON: `{
				"initializationOptions": {
					"test": true,
					"count": 42
				}
			}`,
			validate: func(t *testing.T, config *ServerConfig) {
				if config.InitializationOptions == nil {
					t.Fatal("expected non-nil InitializationOptions")
				}
				if val, ok := config.InitializationOptions["test"].(bool); !ok || !val {
					t.Error("expected test=true")
				}
				if val, ok := config.InitializationOptions["count"].(float64); !ok || val != 42 {
					t.Error("expected count=42")
				}
			},
		},
		{
			name: "valid config with settings",
			configJSON: `{
				"settings": {
					"python": {
						"analysis": {
							"typeCheckingMode": "strict"
						}
					}
				}
			}`,
			validate: func(t *testing.T, config *ServerConfig) {
				if config.Settings == nil {
					t.Fatal("expected non-nil Settings")
				}
				pythonSettings, ok := config.Settings["python"].(map[string]any)
				if !ok {
					t.Fatal("expected python settings")
				}
				analysisSettings, ok := pythonSettings["analysis"].(map[string]any)
				if !ok {
					t.Fatal("expected analysis settings")
				}
				if mode, ok := analysisSettings["typeCheckingMode"].(string); !ok || mode != "strict" {
					t.Error("expected typeCheckingMode=strict")
				}
			},
		},
		{
			name:        "invalid JSON",
			configJSON:  `{invalid json}`,
			expectError: true,
		},
		{
			name:        "file not found",
			configJSON:  "use-file-path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var configPath string

			if tt.configJSON == "use-file-path" {
				configPath = "/non/existent/file.json"
			} else if tt.configJSON != "" {
				// Create temp file with config
				tmpDir := t.TempDir()
				configPath = filepath.Join(tmpDir, "config.json")
				if err := os.WriteFile(configPath, []byte(tt.configJSON), 0644); err != nil {
					t.Fatal(err)
				}
			}

			config, err := LoadServerConfig(configPath)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	tests := []struct {
		command  string
		validate func(t *testing.T, config *ServerConfig)
	}{
		{
			command: "gopls",
			validate: func(t *testing.T, config *ServerConfig) {
				if config.InitializationOptions == nil {
					t.Fatal("expected non-nil InitializationOptions for gopls")
				}
				codelenses, ok := config.InitializationOptions["codelenses"].(map[string]bool)
				if !ok {
					t.Fatal("expected codelenses configuration")
				}
				if !codelenses["test"] {
					t.Error("expected test codelens to be enabled")
				}
			},
		},
		{
			command: "pyright",
			validate: func(t *testing.T, config *ServerConfig) {
				if config.InitializationOptions != nil {
					t.Error("expected nil InitializationOptions for unknown LSP")
				}
				if config.Settings != nil {
					t.Error("expected nil Settings for unknown LSP")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			config := GetDefaultConfig(tt.command)
			if config == nil {
				t.Fatal("expected non-nil config")
			}
			tt.validate(t, config)
		})
	}
}
