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

- Starts a **gRPC server**.
- Registers a `Sync` method, which implements the expected interface. The main app will send options to the connector. The connector developers should give the users the documentation of their connector and the required options that the user should provide and in which format. The app will send a raw string to the connector.
- Uses the callback function to **send paginated data** back to the main application.

The interface definition for connectors can be found in the **`pkg/connector`** folder.

For detailed examples, refer to:

- The [`pkg/connector`](.) directory in this repository.
- The official **[HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)** documentation.
