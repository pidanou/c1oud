# c1: Free, Extensible, Open-Source Data Aggregator

![Screenshot](docs/screenshot.png)

c1 is a free, extensible, open-source data aggregator designed to list all files and data from any location. It allows users to integrate various cloud storage services into a single interface, providing an easy way to find your items.

## Features

- **Extensible**: Easily add support for new storage services with connectors.
- **Open-Source**: Community-driven development and transparency.
- **Cross-Platform**: Runs on multiple operating systems.
- **User-Friendly**: Simple interface for managing files across different storage services.

## How to use

- Install connectors
- Create accounts and configure them
- Sync your account
- Done

## WIP

- CLI version
- CRON jobs
- Metadata of connector syncs to allow differential sync
- Connector cleanup/reinstall
- Find a system to allow deletion of deleted items

## Getting Started

### Prerequisites

- Go 1.24 or later
- Make
- [templ](https://templ.guide/)

### Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/c1-core.git
   cd c1-core
   ```

2. **Download Dependencies**:
   ```bash
   go mod download
   ```

### Set the environment variables

```
# Only needed for Postgres
C1_POSTGRES_DSN="postgres://postgres:password@localhost:5432/postgres?sslmode=disable

# Defaults to no env which serves the static file from embedded fs and will not work on local dev
ENV="dev"

# Either sqlite or postgres defaults to sqlite
C1_DB_ENGINE=postgres

# Defaults to :7777
PORT=":7777"
```

### Running the application

To run c1, use the following command:

```bash
make dev
```

This command will start the application in live mode. The app will be available at `localhost:1323`

### Building the Application

To build the application for deployment, use:

```bash
make build-webapp
```

This will generate the binary file `c1` in the project root directory.

## Makefile Targets

- **`dev`**: Runs the web application in live mode.
- **`dev-templ`**: Runs the template generator in watch mode.
- **`build-webapp`**: Generates templates and builds the web application.

## Extend c1 with connectors

c1 is designed to be extensible, allowing you to easily add support for new storage services by developing connectors.

### Create a connector

See [`pkg/connector`](./pkg/connector)

### Install

A connector is installed by reading a JSON. The JSON can be fetched from a remote source, a local file, or raw JSON string. Here is an example of JSON config. The connector name must be unique, you can override the plugin name even if you fetch it from a remote source.

```json
{
  "name": "s3",
  "description": "A simple connector for S3",
  "source": "VCS",
  "uri": "https://github.com/pidanou/c1-core",
  "install_command": "go build -o s3 pkg/connector/s3/s3.go && chmod +x s3",
  "update_command": "",
  "command": "./s3/s3"
}
```

## Contributing

We welcome contributions from the community! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive messages.
4. Open a pull request.

## License

c1 is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For questions or support, please open an issue on the GitHub repository.

---
