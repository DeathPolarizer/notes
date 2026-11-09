package main

import (
	"context"
	"flag"
	"log"
	"time"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	title := flag.String("title", "", "note title")
	get := flag.String("get", "", "get note by id")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	client := notesv1.NewNotesServiceClient(conn)

	if *title != "" {
		resp, err := client.CreateNote(ctx, &notesv1.CreateNoteRequest{
			Title:   *title,
			Content: "hello",
		})
		if err != nil {
			log.Fatalf("create: %v", err)
		}
		log.Printf("created: %s", resp.Id)
	}

	if *get != "" {
		resp, err := client.GetNote(ctx, &notesv1.GetNoteRequest{
			Id: *get,
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Fatalf("get: code=%s msg=%s", st.Code(), st.Message())
		}
		log.Printf("Got: %+v", resp.Note)
	}
}
