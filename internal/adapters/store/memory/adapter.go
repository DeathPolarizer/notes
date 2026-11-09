package memory

import (
	"context"
	"sync"

	notesv1 "github.com/DeathPolarizer/notes/gen/notes/v1"
	"github.com/DeathPolarizer/notes/internal/adapters/store"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]*notesv1.Note
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]*notesv1.Note)}
}

func (m *MemoryStore) CreateNote(_ context.Context, title, content string) (*notesv1.Note, error) {
	note := &notesv1.Note{
		Id:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
	m.mu.Lock()
	m.data[note.Id] = note
	m.mu.Unlock()
	return proto.Clone(note).(*notesv1.Note), nil
}

func (m *MemoryStore) GetNote(_ context.Context, id string) (*notesv1.Note, error) {
	m.mu.RLock()
	note, ok := m.data[id]
	m.mu.RUnlock()
	if !ok {
		return nil, store.ErrStoreNotFound
	}
	return note, nil
}
