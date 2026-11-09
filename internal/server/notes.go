package server

import (
	"context"
	"errors"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
	"github.com/DeathPolarizer/notes/internal/adapters/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotesServer struct {
	notesv1.UnimplementedNotesServiceServer
	store store.Store
}

func NewNotesServer(store store.Store) *NotesServer {
	return &NotesServer{
		store: store,
	}
}

func (s *NotesServer) CreateNote(ctx context.Context, req *notesv1.CreateNoteRequest) (*notesv1.CreateNoteResponse, error) {
	note, err := s.store.CreateNote(ctx, req.GetTitle(), req.GetContent())
	if err != nil {

	}
	return &notesv1.CreateNoteResponse{Id: note.Id}, nil
}

func (s *NotesServer) GetNote(ctx context.Context, req *notesv1.GetNoteRequest) (*notesv1.GetNoteResponse, error) {
	note, err := s.store.GetNote(ctx, req.GetId())
	if errors.Is(err, store.ErrStoreNotFound) {
		return nil, status.Error(codes.NotFound, err.Error()) // TODO: add constant err message
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error()) // TODO: add constant err message
	}
	return &notesv1.GetNoteResponse{Note: note}, nil
}

// func (UnimplementedNotesServiceServer) ListNotes(context.Context, *ListNotesRequest) (*ListNotesResponse, error) {
// 	return nil, status.Error(codes.Unimplemented, "method ListNotes not implemented")
// }
// func (UnimplementedNotesServiceServer) UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error) {
// 	return nil, status.Error(codes.Unimplemented, "method UpdateNote not implemented")
// }
// func (UnimplementedNotesServiceServer) DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error) {
// 	return nil, status.Error(codes.Unimplemented, "method DeleteNote not implemented")
// }
