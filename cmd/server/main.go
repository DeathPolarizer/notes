package main

import (
	"log"
	"net"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
	"github.com/DeathPolarizer/notes/internal/adapters/store/memory"
	"github.com/DeathPolarizer/notes/internal/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	notesv1.RegisterNotesServiceServer(srv, server.NewNotesServer(memory.NewMemoryStore()))

	log.Println("listening on port :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
