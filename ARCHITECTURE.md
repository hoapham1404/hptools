# HP Tools - Refactored Architecture

## Overview

HP Tools has been refactored from a monolithic structure to a clean, maintainable architecture following Go best practices and design patterns.

## Project Structure

```
hptools/
├── cmd/hptools/           # Application entry points (currently unused)
├── internal/              # Private application packages
│   ├── config/           # Configuration management
│   ├── errors/           # Structured error types
│   ├── logging/          # Logging setup and utilities  
│   ├── models/           # Data structures and DTOs
│   ├── services/         # Business logic layer
│   └── windows/          # Windows API wrapper
├── frontend/             # Wails frontend (React/TypeScript)
├── build/               # Build configuration and assets
└── bin/                 # Compiled binaries
```

## Architecture Principles

### 1. Separation of Concerns
- **Models**: Data structures (`ProcessInfo`, `WindowInfo`, etc.)
- **Services**: Business logic (process management, window operations)
- **Windows**: Windows API abstraction layer
- **Config**: Application configuration and defaults
- **Logging**: Structured logging throughout the application

### 2. Dependency Injection
Services are created with their dependencies injected, making the code testable and modular:

```go
api := windows.NewAPI()
windowService := services.NewWindowService(api, logger)
```

### 3. Interface-Based Design
Services implement interfaces for better testability and flexibility:

```go
type ProcessManager interface {
    GetApplicationProcesses() ([]models.ProcessInfo, error)
    // ...
}

type WindowManager interface {
    SetWindowSize(pid int, width, height int) error
    // ...
}
```

## Key Components

### Configuration (`internal/config`)
- Centralized configuration management
- JSON-based config files with sensible defaults
- Environment-specific settings for app, window, and logging

### Services (`internal/services`)
- **ProcessManager**: Handles process discovery and filtering
- **WindowManager**: Manages window positioning and sizing
- **WindowService**: Combined interface for both managers
- **WailsWindowService**: Wails-specific wrapper for frontend binding

### Windows API (`internal/windows`)
- Clean abstraction over Windows system calls
- Testable interface for API operations
- Centralized Windows-specific functionality

### Error Handling (`internal/errors`)
- Structured error types with context
- Error categorization (Process, Window, API, Config)
- Proper error wrapping and unwrapping

### Logging (`internal/logging`)
- Structured logging with `slog`
- Configurable log levels and formats
- Component-based logger creation

## Benefits of This Architecture

### 1. **Maintainability**
- Clear separation of responsibilities
- Easy to locate and modify specific functionality
- Consistent patterns throughout the codebase

### 2. **Testability**
- Interface-based design enables easy mocking
- Dependency injection allows isolated unit testing
- Clear boundaries between components

### 3. **Extensibility**
- New features can be added as separate services
- Configuration system supports new settings
- Plugin-like architecture for future enhancements

### 4. **Debugging**
- Structured logging with context
- Typed errors with detailed information
- Clear error propagation paths

## Usage Examples

### Adding a New Service

1. Define the interface in `internal/services/interfaces.go`
2. Implement the service in a new file
3. Register it in the main application
4. Create a Wails wrapper if needed for frontend binding

### Adding Configuration Options

1. Update the config structs in `internal/config/config.go`
2. Update the default configuration
3. Use the new settings in your services

### Error Handling

```go
if err != nil {
    return errors.NewWindowError("failed to resize window", err)
}
```

## Future Enhancements

The new architecture makes it easy to add:

- **Plugin System**: Services can be loaded dynamically
- **Different Backends**: Alternative implementations of managers
- **Configuration Hot-reload**: Watch config files for changes
- **Metrics and Monitoring**: Add observability components
- **Testing Framework**: Comprehensive unit and integration tests

## Migration Notes

- All frontend bindings remain the same - no frontend changes needed
- Configuration is now externalized and can be customized
- Logging is more detailed and structured
- Error messages are more informative

This refactored architecture provides a solid foundation for future development while maintaining all existing functionality.