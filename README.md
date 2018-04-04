gRPC prototype
===

Multi-language gRPC prototype.

## Javascript Client

To generate client's `polygon_pb.js`

```shell
protoc --proto_path=definitions/ --js_out=import_style=commonjs,binary:client/ definitions/polygon.proto
```

## golang Server

To generate gRPC and golang structs for the server side

```shell
protoc --go_out=./server definitions/polygon.proto
protoc -I definitions/ --go_out=plugins=grpc:server definitions/polygon.proto
```
