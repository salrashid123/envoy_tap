package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/http"
	//envoy_config_tap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/tap/v3"
	//envoy_config_tap_v3_pb "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
)

var (
	address = flag.String("host", "localhost:9000", "host:port of gRPC server")
)

func init() {
}

func main() {
	flag.Parse()

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
	// var prettyJSON bytes.Buffer
	// err = json.Indent(&prettyJSON, body, "", "\t")
	// if err != nil {
	// 	log.Println("JSON parse error: ", err)
	// 	return
	// }
	// log.Printf("%s\n", string(prettyJSON.Bytes()))

	resp, err := http.Post("http://localhost:9000/tap", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Status: [%s]\n", resp.Status)
	fmt.Println()
	reader := bufio.NewReader(resp.Body)
	for {
		// just print the responses:
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", string(line))

		// streaming: false
		// unmashall as envoy_config_tap_v3_pb.HttpBufferedTrace{}

		// if streaming: true
		// unmarshal as envoy_config_tap_v3_pb.HttpStreamedTraceSegment{}
	}

}
