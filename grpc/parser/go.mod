module main

go 1.17

require (
	github.com/envoyproxy/go-control-plane v0.9.10-0.20210907150352-cf90f659a021
	github.com/golang/protobuf v1.5.0
	github.com/robinpowered/go-proto v0.0.0
	github.com/salrashid123/envoy_tap/echo v0.0.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/cncf/xds/go v0.0.0-20211011173535-cb28da3451f1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.42.0 // indirect
)

replace github.com/salrashid123/envoy_tap/echo => ../src/echo

