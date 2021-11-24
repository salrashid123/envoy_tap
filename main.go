package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	//envoy_config_tap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/tap/v3"
	//envoy_config_tap_v3_pb "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
	envoy_config_tap_v3_pb "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
	//"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/jsonpb"
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
    max_buffered_rx_bytes: 5000
    max_buffered_tx_bytes: 5000		
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
	fmt.Println()
	var rb []byte
	reader := bufio.NewReader(resp.Body)
	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading streamed bytes %v", err)
		}
		fmt.Printf("%s", string(line))
		rb = append(rb, line...)
		// streaming: false
		// unmashall as envoy_config_tap_v3_pb.HttpBufferedTrace{}
		var tw envoy_config_tap_v3_pb.TraceWrapper
		rdr := bytes.NewReader(rb)

		// just wait until you get a proper TraceWrapper
		//  i don't really like this approach since rb can be unbounded string..
		err = jsonpb.Unmarshal(rdr, &tw)
		if err == nil {
			bt := tw.GetHttpBufferedTrace()
			pbody := bt.Response.Body.GetAsBytes()
			log.Printf("Message %s\n", string(pbody))
		}
		// if streaming: true
		// unmarshal as envoy_config_tap_v3_pb.HttpStreamedTraceSegment{}

	}

}
