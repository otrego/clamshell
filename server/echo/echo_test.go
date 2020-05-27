// This is an integration-style test of the echo server.
package echo

import (
	"context"
	"log"
	"net"
	"strings"
	"testing"

	pb "github.com/otrego/clamshell/server/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterEchoServiceServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)
	resp, err := client.GetEchoMessage(ctx, &pb.EchoRequest{
		Id: "foo",
	})
	if err != nil {
		t.Fatalf("GetEchoMessage failed: %v", err)
	}
	if !strings.Contains(resp.GetContent(), "foo") {
		t.Errorf("Response got %v, but expected it to contain \"foo\"", resp)
	}
}
