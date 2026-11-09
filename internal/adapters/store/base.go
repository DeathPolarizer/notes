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
	// GetNotes(context.Context, int32, string) ([]*notesv1.Note, string)
	// UpdateNote(context.Context, string, string, string) *notesv1.Note
	// DeleteNote(context.Context, string)
}
