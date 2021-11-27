## Sample gRPC Client/Server

1. get envoy 
  docker cp `docker create envoyproxy/envoy-dev:latest`:/usr/local/bin/envoy .

2. optional compile protos

$ protoc --version
libprotoc 3.19.1

protoc -I/usr/local/include -I . --go_out=. \
    --descriptor_set_out=src/echo/echo.pb --go_opt=paths=source_relative \
    --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative src/echo/echo.proto

3. run server
  go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --tlsCert=certs/grpc_server_crt.pem \
    --tlsKey=certs/grpc_server_key.pem


4. run envoy

./envoy -c basic.yaml


5. run tap client


 go run parser/main.go


6. run grpc_client

  go run src/grpc_client.go \
    --host 127.0.0.1:8080 \
    --tlsCert=certs/CA_crt.pem \
    --servername=grpc.domain.com



https://github.com/salrashid123/grpc_curl



https://github.com/golang/protobuf/issues/1340

descriptorpb.DescriptorProto describes a protobuf message type. A DescriptorProto is itself a protocol buffer, and protoc can produce these from a .proto file.

protoreflect.MessageDescriptor also describes a protobuf message type. This is a more convenient representation for Go programs than the raw DescriptorProto. The protodesc package can convert between a DescriptorProto and a MessageDescriptor.

protogen.Message is a MessageDescriptor with some additional information that is useful for code generators.

proto.Message is an instance of a message. All the types above describe a type of message; proto.Message is an instantiation of that type.
The dynamicpb package provides a way to take an arbitrary protoreflect.MessageDescriptor and create a proto.Message.