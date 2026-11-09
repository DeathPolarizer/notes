package store

import (
	"context"
	"errors"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
)

var (
	ErrStoreNotFound = errors.New("note not found")
)

type Store interface {
	CreateNote(ctx context.Context, title, content string) (*notesv1.Note, error)
	GetNote(ctx context.Context, id string) (*notesv1.Note, error)
	ListNotes(ctx context.Context, pageSize int32, offset int32) ([]*notesv1.Note, int32)
	UpdateNote(ctx context.Context, note *notesv1.Note) (*notesv1.Note, error)
	DeleteNote(ctx context.Context, id string)
}
