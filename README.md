# hptools

hptools is a desktop application built with Wails3, designed to provide a suite of productivity and utility tools for developers and power users.

## Features

- Modular toolset for various developer workflows
- Modern desktop UI with fast performance
- Cross-platform support (Windows, macOS, Linux)
- Easy integration with custom scripts and plugins

## Getting Started

1. **Install dependencies** (ensure you have Go and Node.js installed).

2. **Development mode**  
   Run the app with hot-reloading:
   ```
   wails3 dev
   ```

3. **Production build**  
   Build a standalone executable:
   ```
   wails3 build
   ```
   The output will be in the `build` directory.

## Project Structure

- `frontend/` — Frontend code (HTML, CSS, JS/TS)
- `main.go` — Go backend entry point
- `app.go` — Application logic and exported methods
- `wails.json` — Wails project configuration

## Documentation

- [Wails3 Documentation](https://v3alpha.wails.io/)

## Contributing

Pull requests and issues are welcome! Please open an issue to discuss your ideas or report bugs.

## License

[MIT](LICENSE)

---

Happy coding with hptools!
