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
	"github.com/psanford/lencode"

	"github.com/salrashid123/envoy_tap/echo"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/protobuf/proto"
	//pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
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
    - name: ":path"
      exact_match: "/echo.EchoServer/SayHelloUnary"
  output_config:
    streaming: false
    max_buffered_rx_bytes: 5000
    max_buffered_tx_bytes: 5000		
    sinks:
    - format: JSON_BODY_AS_BYTES
      streaming_admin: {}`

	// - name: ":path"
	//   exact_match: "/echo.EchoServer/SayHelloUnary"
	// - name: ":path"
	//   exact_match: "/google.pubsub.v1.Publisher/ListTopics"

	body := []byte(c)

	resp, err := http.Post("http://localhost:9000/tap", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error connecting to envoy admin port %v\n", err)
	}
	defer resp.Body.Close()
	log.Printf("Status: [%s]\n", resp.Status)
	log.Println()

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
		rb = append(rb, line...)
		// streaming: false
		// unmashall as envoy_config_tap_v3_pb.HttpBufferedTrace{}
		// if streaming: true
		// unmarshal as envoy_config_tap_v3_pb.HttpStreamedTraceSegment{}
		var tw envoy_config_tap_v3_pb.TraceWrapper
		rdr := bytes.NewReader(rb)
		err = jsonpb.Unmarshal(rdr, &tw)
		if err == nil {
			bt := tw.GetHttpBufferedTrace()
			in := bt.Response.Body.GetAsBytes()
			buf := bytes.NewBuffer(in)
			dec := lencode.NewDecoder(buf, lencode.SeparatorOpt([]byte{0}))
			for {
				msg, err := dec.Decode()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatalf("could not Decode  %v", err)
				}
				pm := echo.EchoReply{}
				//pm := pubsubpb.ListTopicsResponse{}
				err = proto.Unmarshal(msg, &pm)
				if err != nil {
					log.Fatalf("could not Unmarshall  %v", err)
				}
				fmt.Printf("%v\n", pm)
			}
			rb = []byte("")
		}

	}

}
