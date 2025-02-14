# C1oud: Free, Extensible, Open-Source Storage Aggregator

C1oud is a free, extensible, open-source storage aggregator designed to list all files of multiple storage solutions. It allows users to integrate various cloud storage services into a single interface, providing an easy way to find your items.

## Features

- **Extensible**: Easily add support for new storage services with plugins.
- **Open-Source**: Community-driven development and transparency.
- **Cross-Platform**: Runs on multiple operating systems.
- **User-Friendly**: Simple interface for managing files across different storage services.

## Getting Started

### Prerequisites

- Go 1.18 or later
- Make

### Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/c1oud.git
   cd c1oud
   ```

2. **Download Dependencies**:
   ```bash
   go mod download
   ```

### Running the Application

To run C1oud, use the following command:

```bash
make run
```

This command will start the application in live mode, allowing you to manage your storage services interactively.

### Building the Application

To build the application for deployment, use:

```bash
make build-webapp
```

This will generate the binary file `c1oud` in the project root directory.

## Makefile Targets

- **`run-webapp`**: Runs the web application in live mode.
- **`run-templ`**: Runs the template generator in watch mode.
- **`run`**: Runs both `run-templ` and `run-webapp` concurrently.
- **`build-webapp`**: Generates templates and builds the web application.

## Extend C1oud with plugins

C1oud is designed to be extensible, allowing you to easily add support for new storage services by following these steps:

[TODO]

## Contributing

We welcome contributions from the community! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive messages.
4. Open a pull request.

## License

C1oud is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For questions or support, please open an issue on the GitHub repository.

---
