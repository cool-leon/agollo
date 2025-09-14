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

package agollo_test

import (
	"fmt"
	"log"

	"github.com/apolloconfig/agollo/v4"
)

// This example demonstrates how to use the new ApolloWrapper
// to replace the original initialization code from the problem statement
func ExampleNewApolloWrapper() {
	// Original code that needed improvement:
	/*
		func InitApollo() {
			c := &config.AppConfig{
				AppID:          "testApplication_yang",
				Cluster:        "dev",
				IP:             "http://49.235.66.235:8080",
				NamespaceName:  "application",
				IsBackupConfig: true,
				Secret:         "04d35b5b9d264e948e8b8364e5ab62b6",
			}

			client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
				return c, nil
			})

			ApolloClient = &client

			fmt.Println("初始化Apollo配置成功")

			cache := client.GetConfigCache(c.NamespaceName)
			value, _ := cache.Get("/health")
			fmt.Println(value)
		}
	*/

	// New improved code using ApolloWrapper:
	cfg := &agollo.ClientConfig{
		AppID:          "testApplication_yang",
		Cluster:        "dev",
		IP:             "http://49.235.66.235:8080",
		NamespaceNames: []string{"application"},
		IsBackupConfig: true,
		Secret:         "04d35b5b9d264e948e8b8364e5ab62b6",
	}

	// Initialize Apollo client with proper error handling
	apolloClient, err := agollo.NewApolloWrapper(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Apollo client: %v", err)
	}
	defer apolloClient.Close()

	fmt.Println("初始化Apollo配置成功")

	// Get configuration values with better error handling and type safety
	healthValue, err := apolloClient.GetString("/health")
	if err != nil {
		log.Printf("Failed to get /health config: %v", err)
	} else {
		fmt.Printf("Health status: %s\n", healthValue)
	}

	// Examples of type-safe configuration retrieval:

	// Get string with default value
	dbHost := apolloClient.GetStringWithDefault("db.host", "localhost")
	fmt.Printf("Database host: %s\n", dbHost)

	// Get integer configuration
	dbPort := apolloClient.GetIntWithDefault("db.port", 3306)
	fmt.Printf("Database port: %d\n", dbPort)

	// Get boolean configuration
	enableDebug := apolloClient.GetBoolWithDefault("debug.enabled", false)
	fmt.Printf("Debug enabled: %t\n", enableDebug)

	// Get float configuration
	timeout := apolloClient.GetFloatWithDefault("request.timeout", 30.0)
	fmt.Printf("Request timeout: %.1f seconds\n", timeout)

	// Get string slice configuration
	allowedHosts := apolloClient.GetStringSliceWithDefault("security.allowed_hosts", []string{"localhost"})
	fmt.Printf("Allowed hosts: %v\n", allowedHosts)

	// Working with multiple namespaces
	if len(apolloClient.GetNamespaces()) > 1 {
		// Get configuration from specific namespace
		customValue, err := apolloClient.GetStringFromNamespace("custom.setting", "custom")
		if err != nil {
			log.Printf("Failed to get custom setting: %v", err)
		} else {
			fmt.Printf("Custom setting: %s\n", customValue)
		}
	}

	// Access underlying client for advanced operations if needed
	underlyingClient := apolloClient.GetClient()
	_ = underlyingClient // Use for advanced operations
}

// Example of multiple namespace configuration
func ExampleNewApolloWrapper_multipleNamespaces() {
	cfg := &agollo.ClientConfig{
		AppID:          "myApp",
		Cluster:        "production",
		IP:             "http://apollo-server:8080",
		NamespaceNames: []string{"application", "database", "cache"},
		IsBackupConfig: true,
	}

	apolloClient, err := agollo.NewApolloWrapper(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Apollo client: %v", err)
	}
	defer apolloClient.Close()

	// Get values from different namespaces
	appName := apolloClient.GetStringWithDefaultFromNamespace("app.name", "MyApp", "application")
	dbURL := apolloClient.GetStringWithDefaultFromNamespace("url", "localhost:5432", "database")
	cacheSize := apolloClient.GetIntWithDefaultFromNamespace("max_size", 1000, "cache")

	fmt.Printf("App: %s, DB: %s, Cache Size: %d\n", appName, dbURL, cacheSize)
}

// Example of error handling
func ExampleNewApolloWrapper_errorHandling() {
	// Invalid configuration - missing required fields
	cfg := &agollo.ClientConfig{
		// Missing AppID, Cluster, and IP
		NamespaceNames: []string{"application"},
	}

	apolloClient, err := agollo.NewApolloWrapper(cfg)
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		return
	}
	defer apolloClient.Close()

	// This won't be reached due to validation error above
	value, err := apolloClient.GetString("some.key")
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
	} else {
		fmt.Printf("Value: %s\n", value)
	}
}