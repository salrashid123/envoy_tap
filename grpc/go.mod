module main

go 1.13

require (
	github.com/envoyproxy/go-control-plane v0.9.10-0.20210907150352-cf90f659a021
	github.com/golang/protobuf v1.5.2
	github.com/psanford/lencode v0.3.0 // indirect
	github.com/salrashid123/envoy_tap/echo v0.0.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/salrashid123/envoy_tap/echo => ./src/echo
