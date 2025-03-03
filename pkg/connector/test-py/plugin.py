# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

from concurrent import futures
import sys
import time

import grpc

from proto import connector_pb2
from proto import connector_pb2_grpc

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc

class ConnectorServicer(connector_pb2_grpc.ConnectorServicer):
    def Sync(self, request, context):
        print(request.options)
        # write_to_file("/Users/pidanou/.c1/tmp/log", "sd")
        result = connector_pb2.EndSync()
        # broker_address = ""
        # channel = grpc.insecure_channel(broker_address)
        # callback_client = connector_pb2_grpc.CallbackHandlerStub(channel)
        # response_data = connector_pb2.SyncResponse()
        result.metadata = str(request)
        return result
        # response=[
        #         connector_pb2.DataObject(
        #             remote_id="123",
        #             resource_name="test_resource",
        #             uri="http://example.com",
        #             metadata="test metadata"
        #         ),
        #         connector_pb2.DataObject(
        #             remote_id="456",
        #             resource_name="another_resource",
        #             uri="http://example.com/another",
        #             metadata="more metadata"
        #         )
        #     ]
        # try:
        #     err = callback_client.Callback(response_data)
        #     print("Callback successfully sent to Go plugin.")
        #     result.metadata="ok"
        #     return result
        # except grpc.RpcError as e:
        #     print(f"Failed to send callback: {e}")
        #     result.metadata=e.details()+e.debug_error_string()
        #     return result


def serve():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    connector_pb2_grpc.add_ConnectorServicer_to_server(ConnectorServicer(), server)
    health_pb2_grpc.add_HealthServicer_to_server(health, server)
    server.add_insecure_port('127.0.0.1:1234')
    server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:1234|grpc")
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()


def write_to_file(filename: str, data: str, mode: str = "w"):
    """
    Writes data to a file.

    :param filename: Name of the file to write to.
    :param data: Data to be written to the file.
    :param mode: File open mode ('w' for write, 'a' for append, 'wb' for binary write).
    """
    try:
        with open(filename, mode) as file:
            file.write(data)
        print(f"Successfully wrote to {filename}")
    except Exception as e:
        print(f"Error writing to file: {e}")


