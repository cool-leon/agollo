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
	"net/http"
	"testing"

	. "github.com/tevid/gohamcrest"
	"github.com/apolloconfig/agollo/v4/env/config"
)

func TestNewApolloWrapper_NilConfig(t *testing.T) {
	wrapper, err := NewApolloWrapper(nil)
	Assert(t, wrapper, NilVal())
	Assert(t, err, NotNilVal())
	Assert(t, err.Error(), StartWith("client configuration cannot be nil"))
}

func TestNewApolloWrapper_InvalidConfig(t *testing.T) {
	// Test missing AppID
	cfg := &ClientConfig{
		Cluster: "dev",
		IP:      "http://localhost:8080",
	}
	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, wrapper, NilVal())
	Assert(t, err, NotNilVal())
	Assert(t, err.Error(), StartWith("invalid configuration"))

	// Test missing Cluster
	cfg = &ClientConfig{
		AppID: "testApp",
		IP:    "http://localhost:8080",
	}
	wrapper, err = NewApolloWrapper(cfg)
	Assert(t, wrapper, NilVal())
	Assert(t, err, NotNilVal())
	Assert(t, err.Error(), StartWith("invalid configuration"))

	// Test missing IP
	cfg = &ClientConfig{
		AppID:   "testApp",
		Cluster: "dev",
	}
	wrapper, err = NewApolloWrapper(cfg)
	Assert(t, wrapper, NilVal())
	Assert(t, err, NotNilVal())
	Assert(t, err.Error(), StartWith("invalid configuration"))
}

func TestNewApolloWrapper_Success(t *testing.T) {
	// Setup mock server using the same pattern as existing tests
	c := &config.AppConfig{
		AppID:         "testApp",
		Cluster:       "dev",
		NamespaceName: "application",
	}
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, c)
	c.IP = server.URL

	cfg := &ClientConfig{
		AppID:          c.AppID,
		Cluster:        c.Cluster,
		IP:             c.IP,
		NamespaceNames: []string{"application"},
		IsBackupConfig: true,
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	Assert(t, wrapper, NotNilVal())
	Assert(t, wrapper.GetNamespaces(), Equal([]string{"application"}))

	// Clean up
	wrapper.Close()
}

func TestApolloWrapper_GetString(t *testing.T) {
	// Setup mock server using the same pattern as existing tests
	c := &config.AppConfig{
		AppID:         "testApp",
		Cluster:       "dev",
		NamespaceName: "application",
	}
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, c)
	c.IP = server.URL

	cfg := &ClientConfig{
		AppID:          c.AppID,
		Cluster:        c.Cluster,
		IP:             c.IP,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing key
	value, err := wrapper.GetString("key1")
	Assert(t, err, NilVal())
	Assert(t, value, Equal("value1"))

	// Test getting non-existing key
	value, err = wrapper.GetString("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, Equal(""))
}

func TestApolloWrapper_GetStringWithDefault(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing key
	value := wrapper.GetStringWithDefault("key1", "default")
	Assert(t, value, Equal("value1"))

	// Test getting non-existing key returns default
	value = wrapper.GetStringWithDefault("nonexistent", "default")
	Assert(t, value, Equal("default"))
}

func TestApolloWrapper_GetInt(t *testing.T) {
	// Setup mock server using the same pattern as existing tests
	c := &config.AppConfig{
		AppID:         "testApp",
		Cluster:       "dev",
		NamespaceName: "application",
	}
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, c)
	c.IP = server.URL

	cfg := &ClientConfig{
		AppID:          c.AppID,
		Cluster:        c.Cluster,
		IP:             c.IP,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing int key
	value, err := wrapper.GetInt("intValue")
	Assert(t, err, NilVal())
	Assert(t, value, Equal(123))

	// Test getting non-existing key
	value, err = wrapper.GetInt("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, Equal(0))
}

func TestApolloWrapper_GetIntWithDefault(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing int key
	value := wrapper.GetIntWithDefault("intValue", 999)
	Assert(t, value, Equal(123))

	// Test getting non-existing key returns default
	value = wrapper.GetIntWithDefault("nonexistent", 999)
	Assert(t, value, Equal(999))
}

func TestApolloWrapper_GetBool(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing bool key
	value, err := wrapper.GetBool("boolValue")
	Assert(t, err, NilVal())
	Assert(t, value, Equal(true))

	// Test getting non-existing key
	value, err = wrapper.GetBool("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, Equal(false))
}

func TestApolloWrapper_GetBoolWithDefault(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing bool key
	value := wrapper.GetBoolWithDefault("boolValue", false)
	Assert(t, value, Equal(true))

	// Test getting non-existing key returns default
	value = wrapper.GetBoolWithDefault("nonexistent", true)
	Assert(t, value, Equal(true))
}

func TestApolloWrapper_GetFloat(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing float key
	value, err := wrapper.GetFloat("floatValue")
	Assert(t, err, NilVal())
	Assert(t, value, Equal(123.45))

	// Test getting non-existing key
	value, err = wrapper.GetFloat("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, Equal(0.0))
}

func TestApolloWrapper_GetFloatWithDefault(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing float key
	value := wrapper.GetFloatWithDefault("floatValue", 999.99)
	Assert(t, value, Equal(123.45))

	// Test getting non-existing key returns default
	value = wrapper.GetFloatWithDefault("nonexistent", 999.99)
	Assert(t, value, Equal(999.99))
}

func TestApolloWrapper_MultipleNamespaces(t *testing.T) {
	// Setup mock server with multiple namespaces
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 2)
	handlerMap["application"] = onlyNormalConfigResponse
	handlerMap["custom"] = onlyNormalSecondConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application", "custom"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting value from default namespace
	value, err := wrapper.GetStringFromNamespace("key1", "application")
	Assert(t, err, NilVal())
	Assert(t, value, Equal("value1"))

	// Test getting value from custom namespace
	value, err = wrapper.GetStringFromNamespace("key1-1", "custom")
	Assert(t, err, NilVal())
	Assert(t, value, Equal("value1-1"))

	// Test default namespace behavior
	value, err = wrapper.GetString("key1")
	Assert(t, err, NilVal())
	Assert(t, value, Equal("value1")) // Should use first namespace as default
}

func TestApolloWrapper_GetStringSlice(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing string slice key
	value, err := wrapper.GetStringSlice("stringArray")
	Assert(t, err, NilVal())
	Assert(t, value, Equal([]string{"a", "b", "c"}))

	// Test getting non-existing key
	value, err = wrapper.GetStringSlice("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, NilVal())
}

func TestApolloWrapper_GetIntSlice(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test getting existing int slice key
	value, err := wrapper.GetIntSlice("intArray")
	Assert(t, err, NilVal())
	Assert(t, value, Equal([]int{1, 2, 3}))

	// Test getting non-existing key
	value, err = wrapper.GetIntSlice("nonexistent")
	Assert(t, err, NotNilVal())
	Assert(t, value, NilVal())
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	cfg := &ClientConfig{
		AppID:   "testApp",
		Cluster: "dev",
		IP:      "http://localhost:8080",
	}
	err := validateConfig(cfg)
	Assert(t, err, NilVal())
	Assert(t, cfg.NamespaceNames, Equal([]string{"application"})) // Should set default

	// Test with custom namespaces
	cfg = &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             "http://localhost:8080",
		NamespaceNames: []string{"custom1", "custom2"},
	}
	err = validateConfig(cfg)
	Assert(t, err, NilVal())
	Assert(t, cfg.NamespaceNames, Equal([]string{"custom1", "custom2"})) // Should preserve custom
}

func TestApolloWrapper_GetClient(t *testing.T) {
	// Setup mock server
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request), 1)
	handlerMap["application"] = onlyNormalConfigResponse
	server := runMockConfigFilesServer(handlerMap, nil, appConfig)

	cfg := &ClientConfig{
		AppID:          "testApp",
		Cluster:        "dev",
		IP:             server.URL,
		NamespaceNames: []string{"application"},
	}

	wrapper, err := NewApolloWrapper(cfg)
	Assert(t, err, NilVal())
	defer wrapper.Close()

	// Test that we can get the underlying client
	client := wrapper.GetClient()
	Assert(t, client, NotNilVal())

	// Test that the client works
	value := client.GetValue("key1")
	Assert(t, value, Equal("value1"))
}