// Copyright 2025 Apollo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agollo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apolloconfig/agollo/v4/env/config"
)

// ClientConfig holds the configuration parameters for Apollo client initialization
type ClientConfig struct {
	AppID             string
	Cluster           string
	IP                string
	NamespaceNames    []string
	IsBackupConfig    bool
	BackupConfigPath  string
	Secret            string
	Label             string
	SyncServerTimeout int
	MustStart         bool
}

// ApolloWrapper provides a cleaner interface for Apollo client operations
type ApolloWrapper struct {
	client Client
	config *ClientConfig
}

// NewApolloWrapper creates a new Apollo client wrapper with the provided configuration
func NewApolloWrapper(cfg *ClientConfig) (*ApolloWrapper, error) {
	if cfg == nil {
		return nil, errors.New("client configuration cannot be nil")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Convert to internal config format
	appConfig := &config.AppConfig{
		AppID:             cfg.AppID,
		Cluster:           cfg.Cluster,
		IP:                cfg.IP,
		NamespaceName:     strings.Join(cfg.NamespaceNames, ","),
		IsBackupConfig:    cfg.IsBackupConfig,
		BackupConfigPath:  cfg.BackupConfigPath,
		Secret:            cfg.Secret,
		Label:             cfg.Label,
		SyncServerTimeout: cfg.SyncServerTimeout,
		MustStart:         cfg.MustStart,
	}

	// Start Apollo client
	client, err := StartWithConfig(func() (*config.AppConfig, error) {
		return appConfig, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Apollo client: %w", err)
	}

	return &ApolloWrapper{
		client: client,
		config: cfg,
	}, nil
}

// validateConfig validates the client configuration
func validateConfig(cfg *ClientConfig) error {
	if cfg.AppID == "" {
		return errors.New("AppID is required")
	}
	if cfg.Cluster == "" {
		return errors.New("Cluster is required")
	}
	if cfg.IP == "" {
		return errors.New("IP is required")
	}
	if len(cfg.NamespaceNames) == 0 {
		cfg.NamespaceNames = []string{"application"} // default namespace
	}
	return nil
}

// GetString retrieves a string value for the given key from the default namespace
func (w *ApolloWrapper) GetString(key string) (string, error) {
	return w.GetStringFromNamespace(key, "")
}

// GetStringWithDefault retrieves a string value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetStringWithDefault(key, defaultValue string) string {
	return w.GetStringWithDefaultFromNamespace(key, defaultValue, "")
}

// GetStringFromNamespace retrieves a string value for the given key from the specified namespace
func (w *ApolloWrapper) GetStringFromNamespace(key, namespace string) (string, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return "", fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	value := config.GetValue(key)
	if value == "" {
		return "", fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return value, nil
}

// GetStringWithDefaultFromNamespace retrieves a string value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetStringWithDefaultFromNamespace(key, defaultValue, namespace string) string {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetStringValue(key, defaultValue)
}

// GetInt retrieves an integer value for the given key from the default namespace
func (w *ApolloWrapper) GetInt(key string) (int, error) {
	return w.GetIntFromNamespace(key, "")
}

// GetIntWithDefault retrieves an integer value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetIntWithDefault(key string, defaultValue int) int {
	return w.GetIntWithDefaultFromNamespace(key, defaultValue, "")
}

// GetIntFromNamespace retrieves an integer value for the given key from the specified namespace
func (w *ApolloWrapper) GetIntFromNamespace(key, namespace string) (int, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return 0, fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	// Check if key exists first
	value := config.GetValue(key)
	if value == "" {
		return 0, fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return config.GetIntValue(key, 0), nil
}

// GetIntWithDefaultFromNamespace retrieves an integer value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetIntWithDefaultFromNamespace(key string, defaultValue int, namespace string) int {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetIntValue(key, defaultValue)
}

// GetBool retrieves a boolean value for the given key from the default namespace
func (w *ApolloWrapper) GetBool(key string) (bool, error) {
	return w.GetBoolFromNamespace(key, "")
}

// GetBoolWithDefault retrieves a boolean value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetBoolWithDefault(key string, defaultValue bool) bool {
	return w.GetBoolWithDefaultFromNamespace(key, defaultValue, "")
}

// GetBoolFromNamespace retrieves a boolean value for the given key from the specified namespace
func (w *ApolloWrapper) GetBoolFromNamespace(key, namespace string) (bool, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return false, fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	// Check if key exists first
	value := config.GetValue(key)
	if value == "" {
		return false, fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return config.GetBoolValue(key, false), nil
}

// GetBoolWithDefaultFromNamespace retrieves a boolean value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetBoolWithDefaultFromNamespace(key string, defaultValue bool, namespace string) bool {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetBoolValue(key, defaultValue)
}

// GetFloat retrieves a float64 value for the given key from the default namespace
func (w *ApolloWrapper) GetFloat(key string) (float64, error) {
	return w.GetFloatFromNamespace(key, "")
}

// GetFloatWithDefault retrieves a float64 value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetFloatWithDefault(key string, defaultValue float64) float64 {
	return w.GetFloatWithDefaultFromNamespace(key, defaultValue, "")
}

// GetFloatFromNamespace retrieves a float64 value for the given key from the specified namespace
func (w *ApolloWrapper) GetFloatFromNamespace(key, namespace string) (float64, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return 0, fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	// Check if key exists first
	value := config.GetValue(key)
	if value == "" {
		return 0, fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return config.GetFloatValue(key, 0), nil
}

// GetFloatWithDefaultFromNamespace retrieves a float64 value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetFloatWithDefaultFromNamespace(key string, defaultValue float64, namespace string) float64 {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetFloatValue(key, defaultValue)
}

// GetStringSlice retrieves a string slice value for the given key from the default namespace
func (w *ApolloWrapper) GetStringSlice(key string) ([]string, error) {
	return w.GetStringSliceFromNamespace(key, "")
}

// GetStringSliceWithDefault retrieves a string slice value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetStringSliceWithDefault(key string, defaultValue []string) []string {
	return w.GetStringSliceWithDefaultFromNamespace(key, defaultValue, "")
}

// GetStringSliceFromNamespace retrieves a string slice value for the given key from the specified namespace
func (w *ApolloWrapper) GetStringSliceFromNamespace(key, namespace string) ([]string, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return nil, fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	// Check if key exists first
	value := config.GetValue(key)
	if value == "" {
		return nil, fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return config.GetStringSliceValue(key, ",", nil), nil
}

// GetStringSliceWithDefaultFromNamespace retrieves a string slice value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetStringSliceWithDefaultFromNamespace(key string, defaultValue []string, namespace string) []string {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetStringSliceValue(key, ",", defaultValue)
}

// GetIntSlice retrieves an integer slice value for the given key from the default namespace
func (w *ApolloWrapper) GetIntSlice(key string) ([]int, error) {
	return w.GetIntSliceFromNamespace(key, "")
}

// GetIntSliceWithDefault retrieves an integer slice value for the given key with a default value from the default namespace
func (w *ApolloWrapper) GetIntSliceWithDefault(key string, defaultValue []int) []int {
	return w.GetIntSliceWithDefaultFromNamespace(key, defaultValue, "")
}

// GetIntSliceFromNamespace retrieves an integer slice value for the given key from the specified namespace
func (w *ApolloWrapper) GetIntSliceFromNamespace(key, namespace string) ([]int, error) {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return nil, fmt.Errorf("namespace '%s' not found", namespace)
	}
	
	// Check if key exists first
	value := config.GetValue(key)
	if value == "" {
		return nil, fmt.Errorf("key '%s' not found in namespace '%s'", key, namespace)
	}
	
	return config.GetIntSliceValue(key, ",", nil), nil
}

// GetIntSliceWithDefaultFromNamespace retrieves an integer slice value for the given key with a default value from the specified namespace
func (w *ApolloWrapper) GetIntSliceWithDefaultFromNamespace(key string, defaultValue []int, namespace string) []int {
	if namespace == "" {
		namespace = w.getDefaultNamespace()
	}
	
	config := w.client.GetConfig(namespace)
	if config == nil {
		return defaultValue
	}
	
	return config.GetIntSliceValue(key, ",", defaultValue)
}

// GetNamespaces returns all configured namespaces
func (w *ApolloWrapper) GetNamespaces() []string {
	return w.config.NamespaceNames
}

// GetClient returns the underlying Apollo client for advanced operations
func (w *ApolloWrapper) GetClient() Client {
	return w.client
}

// Close stops the underlying Apollo client
func (w *ApolloWrapper) Close() error {
	if w.client != nil {
		w.client.Close()
	}
	return nil
}

// getDefaultNamespace returns the first configured namespace as default
func (w *ApolloWrapper) getDefaultNamespace() string {
	if len(w.config.NamespaceNames) > 0 {
		return w.config.NamespaceNames[0]
	}
	return "application"
}