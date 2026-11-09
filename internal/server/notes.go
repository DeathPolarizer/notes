package server

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
	"github.com/DeathPolarizer/notes/internal/adapters/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func encodePageToken(offset int32) string {
	offsetStr := strconv.Itoa(int(offset))
	return base64.StdEncoding.EncodeToString([]byte(offsetStr))
}

func (s *NotesServer) handleError(err error) error {
	if errors.Is(err, store.ErrStoreNotFound) {
		return status.Error(codes.NotFound, err.Error()) // TODO: add constant err message
	}
	if errors.Is(err, ErrInvalidArgument) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return status.Error(codes.Internal, err.Error()) // TODO: add constant err message
	}
	return nil
}

func (s *NotesServer) CreateNote(ctx context.Context, req *notesv1.CreateNoteRequest) (*notesv1.CreateNoteResponse, error) {
	note, err := s.store.CreateNote(ctx, req.GetTitle(), req.GetContent())
	if statedErr := s.handleError(err); statedErr != nil {
		return nil, statedErr
	}
	return &notesv1.CreateNoteResponse{Id: note.Id}, nil
}

func (s *NotesServer) GetNote(ctx context.Context, req *notesv1.GetNoteRequest) (*notesv1.GetNoteResponse, error) {
	note, err := s.store.GetNote(ctx, req.GetId())
	if statedErr := s.handleError(err); statedErr != nil {
		return nil, statedErr
	}
	return &notesv1.GetNoteResponse{Note: note}, nil
}

func (s *NotesServer) ListNotes(ctx context.Context, req *notesv1.ListNotesRequest) (*notesv1.ListNotesResponse, error) {
	pageSize, offset, err := validateListNotesRequest(req)
	if statedErr := s.handleError(err); statedErr != nil {
		return nil, statedErr
	}
	notes, nextPageOffset := s.store.ListNotes(ctx, pageSize, offset)
	nextPageToken := ""
	if nextPageOffset != 0 {
		nextPageToken = encodePageToken(nextPageOffset)
	}
	return &notesv1.ListNotesResponse{Notes: notes, NextPageToken: nextPageToken}, nil
}

func (s *NotesServer) UpdateNote(ctx context.Context, req *notesv1.UpdateNoteRequest) (*notesv1.UpdateNoteResponse, error) {
	updateNote, err := s.store.GetNote(ctx, req.GetId())
	if statedErr := s.handleError(err); statedErr != nil {
		return nil, statedErr
	}
	if req.Title != nil {
		updateNote.Title = *req.Title
	}
	if req.Content != nil {
		updateNote.Content = *req.Content
	}
	updateNote.UpdatedAt = timestamppb.Now()
	updateNote, err = s.store.UpdateNote(ctx, updateNote)
	if statedErr := s.handleError(err); statedErr != nil {
		return nil, statedErr
	}
	return &notesv1.UpdateNoteResponse{Note: updateNote}, nil
}

func (s *NotesServer) DeleteNote(ctx context.Context, req *notesv1.DeleteNoteRequest) (*notesv1.DeleteNoteResponse, error) {
	s.store.DeleteNote(ctx, req.GetId())
	return &notesv1.DeleteNoteResponse{}, nil
}
