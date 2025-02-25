# Plugin System Overview

This application supports a plugin system that allows external plugins to run as subprocesses and communicate using **gRPC**. The system is built on **HashiCorp's go-plugin** library, which provides a robust framework for managing plugin lifecycles and communication. Developing a plugin is very straithforward and any language can be used.

## How It Works

1. **Starting the Plugin**  
   The application launches the plugin in a separate subprocess.
2. **Sending Options to the Plugin**  
   The main application sends a set of options (configuration) to the plugin.

3. **Fetching Data**  
   The plugin processes the options and fetches the required data.

4. **Streaming Data Back**  
   Instead of sending all data at once, the plugin calls back to the main application using a callback function (passed as an argument to `Sync`). This allows for **pagination** and efficient data handling.

5. **Exiting**  
   Once the plugin has finished its task, it gracefully exits.

## Developing a Plugin

To create a plugin, you need to build an application that:

- Starts a **gRPC server**.
- Registers a `Sync` method, which implements the expected interface. The main app will send options to the plugin. The plugin developers should give the users the documentation of their plugin and the required options that the user should provide and in which format. The app will send a raw string to the plugin.
- Uses the callback function to **send paginated data** back to the main application.

The interface definition for plugins can be found in the **`pkg/plugin`** folder.

For detailed examples, refer to:

- The [`pkg/plugin`](.) directory in this repository.
- The official **[HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)** documentation.
