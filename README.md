## Envoy TAP filter helloworld

Simple implementation of an [Envoy Tap Filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/tap_filter).  Really, thats it.

All that this repo does is shows the "helloworld" of setting up the TAP filter to write request/response to a file and to use the ADMIN interface to dynamically receive the forked metrics.   This repo demonstrates both `HTTP` and `gRPC` message processing

Use this...well, just to understand the basics...i dont' know how its operationalized beyond mesh frameworks like [istio traffic tapping](https://www.envoyproxy.io/docs/envoy/latest/operations/traffic_tapping).


---

to use, you need golang 1.15+ and envoy installed locally

You can get envoy like this:

```bash
docker cp `docker create  envoyproxy/envoy-dev:latest`:/usr/local/bin/envoy .
```

---

## HTTP
### Tap Static


To TAP traffic and write it a file, set `server.yaml` to use


```yaml
          http_filters:
          - name: envoy.filters.http.tap
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.tap.v3.Tap
              common_config:
                # admin_config:
                #   config_id: test_config_id
                static_config:
                  match_config:
                    http_request_headers_match:
                      headers:
                      - name: foo
                        exact_match: bar
                  output_config:
                    streaming: false
                    sinks:
                    - format: JSON_BODY_AS_BYTES
                      file_per_tap:
                        path_prefix: /tmp/                
                
          - name: envoy.filters.http.router
```


Then run envoy

```bash
./envoy -c server.yaml
```

In a new window, send over a request

```bash
curl -v -H "foo: bar" -H "content-type: application/json"  -H "user: sal" -d '{"foo":"bar"}'  http://localhost:8080/post
```

What you'll see in the `/tmp/` folder is a file like this with a random text prefixed by `_` (eg, `/tmp/_473128456949399711.json`)

```json
{
 "http_buffered_trace": {
  "request": {
   "headers": [
    {
     "key": ":authority",
     "value": "localhost:8080"
    },
    {
     "key": ":path",
     "value": "/post"
    },
    {
     "key": ":method",
     "value": "POST"
    },
    {
     "key": ":scheme",
     "value": "http"
    },
    {
     "key": "user-agent",
     "value": "curl/7.74.0"
    },
    {
     "key": "accept",
     "value": "*/*"
    },
    {
     "key": "foo",
     "value": "bar"
    },
    {
     "key": "content-type",
     "value": "application/json"
    },
    {
     "key": "user",
     "value": "sal"
    },
    {
     "key": "content-length",
     "value": "13"
    },
    {
     "key": "x-forwarded-proto",
     "value": "http"
    },
    {
     "key": "x-request-id",
     "value": "230534ca-ebec-402e-83a8-a9702bf8fd78"
    },
    {
     "key": "x-envoy-expected-rq-timeout-ms",
     "value": "15000"
    }
   ],
   "body": {
    "truncated": false,
    "as_bytes": "eyJmb28iOiJiYXIifQ=="
   },
   "trailers": []
  },
  "response": {
   "headers": [
    {
     "key": ":status",
     "value": "200"
    },
    {
     "key": "date",
     "value": "Tue, 23 Nov 2021 21:33:22 GMT"
    },
    {
     "key": "content-type",
     "value": "application/json"
    },
    {
     "key": "content-length",
     "value": "506"
    },
    {
     "key": "server",
     "value": "envoy"
    },
    {
     "key": "access-control-allow-origin",
     "value": "*"
    },
    {
     "key": "access-control-allow-credentials",
     "value": "true"
    },
    {
     "key": "x-envoy-upstream-service-time",
     "value": "26"
    }
   ],
   "body": {
    "truncated": false,
    "as_bytes": "ewogICJhcmdzIjoge30sIAogICJkYXRhIjogIntcImZvb1wiOlwiYmFyXCJ9IiwgCiAgImZpbGVzIjoge30sIAogICJmb3JtIjoge30sIAogICJoZWFkZXJzIjogewogICAgIkFjY2VwdCI6ICIqLyoiLCAKICAgICJDb250ZW50LUxlbmd0aCI6ICIxMyIsIAogICAgIkNvbnRlbnQtVHlwZSI6ICJhcHBsaWNhdGlvbi9qc29uIiwgCiAgICAiRm9vIjogImJhciIsIAogICAgIkhvc3QiOiAibG9jYWxob3N0IiwgCiAgICAiVXNlciI6ICJzYWwiLCAKICAgICJVc2VyLUFnZW50IjogImN1cmwvNy43NC4wIiwgCiAgICAiWC1BbXpuLVRyYWNlLUlkIjogIlJvb3Q9MS02MTlkNWUyMi0xZGMzYTU4OTVkNWM2MzNiNDE3ZDg2ZTIiLCAKICAgICJYLUVudm95LUV4cGVjdGVkLVJxLVRpbWVvdXQtTXMiOiAiMTUwMDAiCiAgfSwgCiAgImpzb24iOiB7CiAgICAiZm9vIjogImJhciIKICB9LCAKICAib3JpZ2luIjogIjcyLjgzLjY3LjE3NCIsIAogICJ1cmwiOiAiaHR0cHM6Ly9sb2NhbGhvc3QvcG9zdCIKfQo="
   },
   "trailers": []
  }
 }
}
```

the response body is b64encoded JSON response from the upstream

If you set `streaming: true`, you'll see streamed segments

```json
{
 "http_streamed_trace_segment": {
  "trace_id": "11638135731790493088",
  "request_headers": {
   "headers": [
    {
     "key": ":authority",
     "value": "localhost:8080"
```

---

### Admin

The following sets up TAP but instead of using the Static config, the parameters are set via a remote application and streamed back to that app:


On `server.yaml`, first enable the admin config
```yaml
          http_filters:
          - name: envoy.filters.http.tap
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.tap.v3.Tap
              common_config:
                admin_config:
                  config_id: test_config_id
```

Then run envoy

```bash
./envoy -c server.yaml
```

Start the remote TAP monitor:

```bash
go run main.go
```

What the tap monitor server does is _remotely_ program and receive TAP responses

programming is done with a simple POST to `/tap` endpoint of envoy

```golang
	c := `
config_id: test_config_id
tap_config:
  match_config:
   http_request_headers_match:
    headers:
    - name: foo
      exact_match: bar
  output_config:
    streaming: false
    sinks:
    - format: JSON_BODY_AS_BYTES
      streaming_admin: {}`

	body := []byte(c)
	resp, err := http.Post("http://localhost:9000/tap", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Status: [%s]\n", resp.Status)
```

What the specs for TAP state is that once you program envoy endpoint, any inbound request that fulfils the TAP specs will get streamed back to the listener.

So, in a new window, send over a request

```bash
curl -v -H "foo: bar" -H "content-type: application/json"  -H "user: sal" -d '{"foo":"bar"}'  http://localhost:8080/post
```

What you should see is the `HttpBufferedTrace` in the response (since we set `streaming: false`)



```json
$ go run main.go 
Status: [200 OK]

{
 "http_buffered_trace": {
  "request": {
   "headers": [
    {
     "key": ":authority",
     "value": "localhost:8080"
    },
    {
     "key": ":path",
     "value": "/post"
    },
    {
     "key": ":method",
     "value": "POST"
    },
    {
     "key": ":scheme",
     "value": "https"
    },
    {
     "key": "user-agent",
     "value": "curl/7.74.0"
    },
    {
     "key": "accept",
     "value": "*/*"
    },
    {
     "key": "foo",
     "value": "bar"
    },
    {
     "key": "content-type",
     "value": "application/json"
    },
    {
     "key": "user",
     "value": "sal"
    },
    {
     "key": "content-length",
     "value": "13"
    },
    {
     "key": "x-forwarded-proto",
     "value": "http"
    },
    {
     "key": "x-request-id",
     "value": "71292137-0566-4962-ace1-eb2d4ad9120b"
    },
    {
     "key": "x-envoy-expected-rq-timeout-ms",
     "value": "15000"
    }
   ],
   "body": {
    "truncated": false,
    "as_bytes": "eyJmb28iOiJiYXIifQ=="
   },
   "trailers": []
  },
  "response": {
   "headers": [
    {
     "key": ":status",
     "value": "200"
    },
    {
     "key": "date",
     "value": "Fri, 04 Jun 2021 19:51:22 GMT"
    },
    {
     "key": "content-type",
     "value": "application/json"
    },
    {
     "key": "content-length",
     "value": "506"
    },
    {
     "key": "server",
     "value": "envoy"
    },
    {
     "key": "access-control-allow-origin",
     "value": "*"
    },
    {
     "key": "access-control-allow-credentials",
     "value": "true"
    },
    {
     "key": "x-envoy-upstream-service-time",
     "value": "19"
    }
   ],
   "body": {
    "truncated": false,
    "as_bytes": "ewogICJhcmdzIjoge30sIAogICJkYXRhIjogIntcImZvb1wiOlwiYmFyXCJ9IiwgCiAgImZpbGVzIjoge30sIAogICJmb3JtIjoge30sIAogICJoZWFkZXJzIjogewogICAgIkFjY2VwdCI6ICIqLyoiLCAKICAgICJDb250ZW50LUxlbmd0aCI6ICIxMyIsIAogICAgIkNvbnRlbnQtVHlwZSI6ICJhcHBsaWNhdGlvbi9qc29uIiwgCiAgICAiRm9vIjogImJhciIsIAogICAgIkhvc3QiOiAibG9jYWxob3N0IiwgCiAgICAiVXNlciI6ICJzYWwiLCAKICAgICJVc2VyLUFnZW50IjogImN1cmwvNy43NC4wIiwgCiAgICAiWC1BbXpuLVRyYWNlLUlkIjogIlJvb3Q9MS02MGJhODQzYS00ZDU4NGExNjMxN2Q2M2RjMzQyOTczOGQiLCAKICAgICJYLUVudm95LUV4cGVjdGVkLVJxLVRpbWVvdXQtTXMiOiAiMTUwMDAiCiAgfSwgCiAgImpzb24iOiB7CiAgICAiZm9vIjogImJhciIKICB9LCAKICAib3JpZ2luIjogIjcyLjgzLjY3LjE3NCIsIAogICJ1cmwiOiAiaHR0cHM6Ly9sb2NhbGhvc3QvcG9zdCIKfQo="
   },
   "trailers": []
  }
 }
}
```

If you want to see the streaming response, set `streaming: true` and the output will be like

```json
{
 "http_buffered_trace": {
  "request": {
   "headers": [
    {
     "key": ":authority",
     "value": "localhost:8080"
    },
    {
     "key": ":path",
     "value": "/post"
    },
    {
     "key": ":method",
     "value": "POST"
    },
    {
     "key": ":scheme",
     "value": "https"

```

One glaring thing i omitted in this repo is how to parse the JSON data back on the admin listener.

It seems  envoy returns a new line delimited pretty-printed JSON back!...i have no idea how to effectively parse that in
golang...there is certainly away...if you know it, LMK

---

## gRPC

For gRPC we will run a simple client->envoy->server just like the http sample above. The admin tap application will connect to envoy and receive a copy of 
a specific, named RPC that traverses envoy.  The tap application will then attempt to decode the protobuf that is contained in that specif call.

As above,

#### get envoy 

```bash
  docker cp `docker create envoyproxy/envoy-dev:latest`:/usr/local/bin/envoy .
```

(optional compile protos)

```bash
$ protoc --version
libprotoc 3.19.1

protoc -I/usr/local/include -I . --go_out=. \
    --descriptor_set_out=src/echo/echo.pb --go_opt=paths=source_relative \
    --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative src/echo/echo.proto
```

#### run server

```bash
  go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --tlsCert=certs/grpc_server_crt.pem \
    --tlsKey=certs/grpc_server_key.pem
```

now run and enable the admin TAP listener

```
./envoy -c basic.yaml
```

#### run tap client

The TAP application will attempt to parse the Unary requests

```golang
	c := `
config_id: test_config_id
tap_config:
  match_config:
   http_request_headers_match:
    headers:
    - name: ":path"
      exact_match: "/echo.EchoServer/SayHelloUnary"
  output_config:
    streaming: false
    max_buffered_rx_bytes: 5000
    max_buffered_tx_bytes: 5000		
    sinks:
    - format: JSON_BODY_AS_BYTES
      streaming_admin: {}`
```

If you want,you can specify other APIs endpoints or wildcard all endpoints.  For the latter, you will have to selectively parse back the messages.  In the case of the attached rpcs, it
doesn't really matter since they are all return `echo.EchoReply`.  Note, this sample does not support compression `Length-Prefixed-Message` in the [grpc protocol](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md))

```bash
go run parser/main.go
```

#### run grpc_client

```bash
go run src/grpc_client.go \
    --host 127.0.0.1:8080 \
    --tlsCert=certs/CA_crt.pem \
    --servername=grpc.domain.com
```

What you should see in the TAP client is the parsed protobuf messages that traversed envoy.

---

thats all folks,

---

## Other links
Other reference envoy samples

- [Envoy WASM and LUA filters for Certificate Bound Tokens](https://github.com/salrashid123/envoy_cert_bound_token)
- [Envoy mTLS](https://github.com/salrashid123/envoy_mtls)
- [Envoy control plane "hello world"](https://github.com/salrashid123/envoy_control)
- [Envoy for Google Cloud Identity Aware Proxy](https://github.com/salrashid123/envoy_iap)
- [Envoy External Authorization server (envoy.ext_authz) with OPA HelloWorld](https://github.com/salrashid123/envoy_external_authz)
- [Envoy RBAC](https://github.com/salrashid123/envoy_rbac)
- [Envoy Global rate limiting helloworld](https://github.com/salrashid123/envoy_ratelimit)
- [Envoy EDS "hello world"](https://github.com/salrashid123/envoy_discovery)
- [Envoy WASM with external gRPC server](https://github.com/salrashid123/envoy_wasm)
- [Redis AUTH and mTLS with Envoy](https://github.com/salrashid123/envoy_redis)

- [gRPC per method observability with envoy, Istio, OpenCensus and GKE](https://github.com/salrashid123/grpc_stats_envoy_istio#envoy)
- [gRPC XDS](https://github.com/salrashid123/grpc_xds)
- [gRPC ALTS](https://github.com/salrashid123/grpc_alts)

- [gRPC curl](https://github.com/salrashid123/grpc_curl)


https://pkg.go.dev/github.com/envoyproxy/go-control-plane/envoy/data/tap/v3

https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/tap/v3/tap.proto#extensions-filters-http-tap-v3-tap
https://www.envoyproxy.io/docs/envoy/latest/api-v3/data/tap/v3/http.proto


```golang
	envoy_config_tap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/tap/v3"
  envoy_config_tap_v3_pb "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
```
