# Apollo Client Wrapper

This document describes the new Apollo client wrapper that provides better encapsulation, error handling, and type-safe configuration retrieval.

## Overview

The `ApolloWrapper` provides a cleaner, more robust interface for working with Apollo configuration compared to the original raw client usage. It addresses several issues with the original implementation:

1. **Better Error Handling**: Proper error propagation and descriptive error messages
2. **Type Safety**: Dedicated methods for different data types (string, int, bool, float, slices)
3. **Configuration Validation**: Validates required fields during initialization
4. **Multiple Namespace Support**: Easy handling of multiple namespaces
5. **Default Values**: Built-in support for default values
6. **Clean API**: Simplified initialization and usage patterns

## Migration from Original Code

### Before (Original Implementation)
```go
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
```

### After (Using ApolloWrapper)
```go
func InitApollo() {
    cfg := &agollo.ClientConfig{
        AppID:          "testApplication_yang",
        Cluster:        "dev",
        IP:             "http://49.235.66.235:8080",
        NamespaceNames: []string{"application"},
        IsBackupConfig: true,
        Secret:         "04d35b5b9d264e948e8b8364e5ab62b6",
    }

    apolloClient, err := agollo.NewApolloWrapper(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize Apollo client: %v", err)
    }
    defer apolloClient.Close()

    fmt.Println("初始化Apollo配置成功")

    healthValue, err := apolloClient.GetString("/health")
    if err != nil {
        log.Printf("Failed to get /health config: %v", err)
    } else {
        fmt.Printf("Health status: %s\n", healthValue)
    }
}
```

## Basic Usage

### Initialization
```go
cfg := &agollo.ClientConfig{
    AppID:          "myApp",
    Cluster:        "production",
    IP:             "http://apollo-server:8080",
    NamespaceNames: []string{"application"},
    IsBackupConfig: true,
    Secret:         "your-secret-key",
}

client, err := agollo.NewApolloWrapper(cfg)
if err != nil {
    log.Fatalf("Failed to initialize Apollo: %v", err)
}
defer client.Close()
```

### Configuration Retrieval

#### Basic Methods (with error handling)
```go
// String values
dbHost, err := client.GetString("database.host")
if err != nil {
    log.Printf("Database host not configured: %v", err)
}

// Integer values
port, err := client.GetInt("server.port")
if err != nil {
    log.Printf("Server port not configured: %v", err)
}

// Boolean values
debug, err := client.GetBool("debug.enabled")
if err != nil {
    log.Printf("Debug flag not configured: %v", err)
}

// Float values
timeout, err := client.GetFloat("request.timeout")
if err != nil {
    log.Printf("Timeout not configured: %v", err)
}
```

#### Methods with Default Values
```go
// String with default
dbHost := client.GetStringWithDefault("database.host", "localhost")

// Integer with default
port := client.GetIntWithDefault("server.port", 8080)

// Boolean with default
debug := client.GetBoolWithDefault("debug.enabled", false)

// Float with default
timeout := client.GetFloatWithDefault("request.timeout", 30.0)

// String slice with default
hosts := client.GetStringSliceWithDefault("allowed.hosts", []string{"localhost"})

// Integer slice with default
ports := client.GetIntSliceWithDefault("allowed.ports", []int{80, 443})
```

### Multiple Namespaces

#### Configuration
```go
cfg := &agollo.ClientConfig{
    AppID:          "myApp",
    Cluster:        "production", 
    IP:             "http://apollo-server:8080",
    NamespaceNames: []string{"application", "database", "cache"},
}

client, err := agollo.NewApolloWrapper(cfg)
// ...
```

#### Namespace-specific Retrieval
```go
// Get from specific namespace
appName, err := client.GetStringFromNamespace("app.name", "application")
dbURL, err := client.GetStringFromNamespace("connection.url", "database")
cacheSize, err := client.GetIntFromNamespace("max.size", "cache")

// With defaults
appName := client.GetStringWithDefaultFromNamespace("app.name", "MyApp", "application")
dbURL := client.GetStringWithDefaultFromNamespace("connection.url", "localhost:5432", "database")
cacheSize := client.GetIntWithDefaultFromNamespace("max.size", 1000, "cache")
```

## Configuration Options

### ClientConfig Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `AppID` | string | Yes | Apollo application ID |
| `Cluster` | string | Yes | Apollo cluster name |
| `IP` | string | Yes | Apollo server URL |
| `NamespaceNames` | []string | No | List of namespaces (defaults to ["application"]) |
| `IsBackupConfig` | bool | No | Enable configuration backup |
| `BackupConfigPath` | string | No | Path for backup files |
| `Secret` | string | No | Secret key for authentication |
| `Label` | string | No | Configuration label |
| `SyncServerTimeout` | int | No | Server sync timeout in seconds |
| `MustStart` | bool | No | Require successful first sync |

## Error Handling

The wrapper provides comprehensive error handling:

### Initialization Errors
```go
client, err := agollo.NewApolloWrapper(cfg)
if err != nil {
    // Handle configuration validation errors
    // or Apollo client startup errors
    log.Fatalf("Apollo initialization failed: %v", err)
}
```

### Configuration Retrieval Errors
```go
value, err := client.GetString("some.key")
if err != nil {
    // Key not found or namespace not available
    log.Printf("Configuration error: %v", err)
    // Use default or alternative logic
}
```

## Advanced Usage

### Access to Underlying Client
```go
// For advanced operations, access the underlying Apollo client
underlyingClient := client.GetClient()

// Use underlying client methods
config := underlyingClient.GetConfig("custom-namespace")
cache := underlyingClient.GetConfigCache("application")
```

### Cleanup
```go
// Always close the client to stop background processes
defer client.Close()

// Or explicitly
err := client.Close()
if err != nil {
    log.Printf("Error closing Apollo client: %v", err)
}
```

## Benefits

1. **Type Safety**: Eliminates runtime type conversion errors
2. **Error Handling**: Proper error propagation instead of silent failures
3. **Validation**: Configuration validation at startup
4. **Simplicity**: Cleaner API for common operations
5. **Flexibility**: Support for multiple namespaces and default values
6. **Maintainability**: Better code organization and documentation
7. **Backward Compatibility**: Can be used alongside existing Apollo client code

## Best Practices

1. **Always handle errors** from initialization and configuration retrieval
2. **Use default values** for non-critical configuration
3. **Close the client** when done to clean up resources
4. **Validate configuration** early in application startup
5. **Group related configuration** in separate namespaces
6. **Use descriptive error messages** in your application logging