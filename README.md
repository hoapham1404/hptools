# HP Tools - Window Management Application

HP Tools is a desktop application built with Wails3 that provides advanced window management functionality for Windows systems. It helps users organize and control application windows with precision.

## Features

- **Process Discovery**: Automatically detect running applications with visible windows
- **Window Positioning**: Set exact window positions and sizes
- **Smart Filtering**: Filter out system processes and show only relevant applications  
- **Structured Logging**: Comprehensive logging for debugging and monitoring
- **Configurable**: JSON-based configuration system with sensible defaults

## Architecture

This project has been refactored for maintainability using clean architecture principles:

```
hptools/
├── internal/
│   ├── config/          # Configuration management
│   ├── errors/          # Structured error handling
│   ├── logging/         # Logging utilities
│   ├── models/          # Data structures
│   ├── services/        # Business logic
│   └── windows/         # Windows API wrapper
├── frontend/            # React/TypeScript UI
└── main.go             # Application entry point
```

For detailed architecture information, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Getting Started

### Prerequisites
- Go 1.24+ 
- Node.js and npm (for frontend development)
- Wails v3

### Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd hptools
   ```

2. **Run in development mode**:
   ```bash
   wails3 dev
   ```

3. **Build for production**:
   ```bash
   wails3 task windows:build PRODUCTION=true
   ```

### Configuration

HP Tools uses a JSON configuration file. Create `~/.config/hptools/config.json` or use the provided example:

```bash
cp config.example.json ~/.config/hptools/config.json
```

Configuration options include:
- **App settings**: Name, description
- **Window settings**: Default size, position, styling
- **Logging**: Level, format (text/json)

## Usage

The application provides a clean interface for:

1. **Viewing running applications** with their process details
2. **Setting window dimensions** with precise width/height control  
3. **Positioning windows** at specific screen coordinates
4. **Real-time window information** including current size and position

## API Reference

### Main Services

- `GetApplicationProcesses()` - Get all applications with visible windows
- `SetWindowSize(pid, width, height)` - Resize a window by process ID
- `SetWindowPosition(pid, x, y, width, height)` - Move and resize window
- `GetWindowInfo(pid)` - Get current window dimensions and position

## Contributing

1. Follow the established architecture patterns
2. Add appropriate logging and error handling
3. Update tests for new functionality
4. Document any new configuration options

## Project Structure Details

- **Models**: Data structures (`ProcessInfo`, `WindowInfo`)
- **Services**: Business logic with dependency injection
- **Windows API**: Clean abstraction over Windows system calls
- **Configuration**: Centralized, file-based settings
- **Logging**: Structured logging with configurable levels

## License

[MIT](LICENSE)

---

Built with ❤️ using Wails v3
