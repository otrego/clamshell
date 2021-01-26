// Package echo contains the echo server.
package echo

import (
	"context"

	pb "github.com/otrego/clamshell/server/api"
)

// Server is the Echo Server implementation.
type Server struct{}

// messages is a map of id to message content
var messages = map[string]string{
	"foo": "this is a foo message",
	"bar": "this is a bar message",
}

// GetEchoMessage gets an echo message.
func (s *Server) GetEchoMessage(ctx context.Context, req *pb.EchoRequest) (*pb.EchoMessage, error) {
	id := req.GetId()
	m, ok := messages[id]
	if !ok {
		return &pb.EchoMessage{Content: "not found"}, nil
	}
	return &pb.EchoMessage{Content: m}, nil
}

// ListEchoMessages lists all features contained within the given bounding Rectangle.
func (s *Server) ListEchoMessages(ctx context.Context, _ *pb.EmptyRequest) (*pb.EchoMessageCollection, error) {
	var out []*pb.EchoMessage
	for _, m := range messages {
		out = append(out, &pb.EchoMessage{Content: m})
	}
	return &pb.EchoMessageCollection{
		EchoMessages: out,
	}, nil
}
