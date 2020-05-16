// Package echo contains the echo server.
package echo

import (
	"context"

	pb "github.com/otrego/clamshell/server/api"
)

type EchoServer struct {
}

// echoMessages is a map of id to message content
var echoMessages = map[string]string{
	"foo": "this is a foo message",
	"bar": "this is a bar message",
}

// GetEchoMessage gets an echo message.
func (s *EchoServer) GetEchoMessage(ctx context.Context, req *pb.EchoRequest) (*pb.EchoMessage, error) {
	id := req.GetId()
	m, ok := echoMessages[id]
	if !ok {
		return &pb.EchoMessage{Content: "not found"}, nil
	}
	return &pb.EchoMessage{Content: m}, nil
}

// ListEchoMessages lists all features contained within the given bounding Rectangle.
func (s *EchoServer) ListEchoMessages(ctx context.Context, _ *pb.EmptyRequest) (*pb.EchoMessageCollection, error) {
	var out []*pb.EchoMessage
	for _, m := range echoMessages {
		out = append(out, &pb.EchoMessage{Content: m})
	}
	return &pb.EchoMessageCollection{
		EchoMessages: out,
	}, nil
}
