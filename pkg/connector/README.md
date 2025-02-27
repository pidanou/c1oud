# Connector System Overview

This application supports a connector system that allows external connectors to run as subprocesses and communicate using **gRPC**. The system is built on **HashiCorp's go-plugin** library, which provides a robust framework for managing connector lifecycles and communication. Developing a connector is very straithforward and any language can be used.

## How It Works

1. **Starting the Connector**  
   The application launches the connector in a separate subprocess.
2. **Sending Options to the Connector**  
   The main application sends a set of options (configuration) to the connector.

3. **Fetching Data**  
   The connector processes the options and fetches the required data.

4. **Streaming Data Back**  
   Instead of sending all data at once, the connector calls back to the main application using a callback function (passed as an argument to `Sync`). This allows for **pagination** and efficient data handling.

5. **Exiting**  
   Once the connector has finished its task, it gracefully exits.

## Developing a Connector

To create a connector, you need to build an application that:

### Overview

- Starts a **gRPC server**.
- Registers a `Sync` method, which implements the expected interface. The main app will send options to the connector. The connector developers should give the users the documentation of their connector and the required options that the user should provide and in which format. The app will send a raw string to the connector.
- Uses the callback function to **send paginated data** back to the main application.

The interface definition for connectors can be found in the **`pkg/connector`** folder.

For detailed examples, refer to:

- The [`pkg/connector`](.) directory in this repository.
- The official **[HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)** documentation.

### Behaviour

A connector can return any data as long as it is in the format of a list of `DataObject` messages. Connectors need to parse the options.
The data returned by the connector will be upserted in the database using the `RemoteID` field.

## Install a connector

The installation of the connector is done through a JSON:

```json
{
  "name": "s3",
  "description": "A simple connector for S3"
  "source": "VCS",
  "uri": "https://github.com/pidanou/c1-core",
  "install_command": "go build -o s3 pkg/connector/s3/s3.go && chmod +x s3",
  "update_command": "",
  "command": "./s3/s3"
}
```

- name (required) : Default connector name that be overriden by the user
- description (optional) : Short description of the plugin
- source (required) : How to fetch the source, can be: HTTP, VCS, Local
- URI (required) : Location of the connector source, should a url, github repo, local path...
- install_command (optional): Command to install the connector after downloading
- update_command (optional) : Command to update the connector
- command (required) : The command to start the connector from the root of the source, for example: python ./path/to/script
