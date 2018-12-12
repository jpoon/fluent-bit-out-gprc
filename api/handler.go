package api

import (
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) RecordEvents(ctx context.Context, in *Record) (*RecordSummary, error) {
	record := in.GetRecord()
	for k, v := range record {
		log.Printf("%s->%s", k, v)
	}
	log.Printf("%s %s %s", in.GetTimestamp(), in.GetTag(), in.GetRecord())
	return &RecordSummary{EventCount: 99}, nil
}
