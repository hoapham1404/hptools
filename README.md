# Welcome to Your New Wails3 Project!

Congratulations on generating your Wails3 application! This README will guide you through the next steps to get your project up and running.

## Getting Started

1. Navigate to your project directory in the terminal.

2. To run your application in development mode, use the following command:

   ```
   wails3 dev
   ```

   This will start your application and enable hot-reloading for both frontend and backend changes.

3. To build your application for production, use:

   ```
   wails3 build
   ```

   This will create a production-ready executable in the `build` directory.

## Exploring Wails3 Features

Now that you have your project set up, it's time to explore the features that Wails3 offers:

1. **Check out the examples**: The best way to learn is by example. Visit the `examples` directory in the `v3/examples` directory to see various sample applications.

2. **Run an example**: To run any of the examples, navigate to the example's directory and use:

   ```
   go run .
   ```

   Note: Some examples may be under development during the alpha phase.

3. **Explore the documentation**: Visit the [Wails3 documentation](https://v3.wails.io/) for in-depth guides and API references.

4. **Join the community**: Have questions or want to share your progress? Join the [Wails Discord](https://discord.gg/JDdSxwjhGf) or visit the [Wails discussions on GitHub](https://github.com/wailsapp/wails/discussions).

## Project Structure

Take a moment to familiarize with project structure:

- `frontend/`: Contains frontend code (I use React + Typescript template)
- `main.go`: The entry point of Go backend
- `app.go`: Define application structure and methods here
- `wails.json`: Configuration file for Wails project

## Next Steps

1. Modify the frontend in the `frontend/` directory to create your desired UI.
2. Add backend functionality in `main.go`.
3. Use `wails3 dev` to see your changes in real-time.
4. When ready, build your application with `wails3 build`.

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
