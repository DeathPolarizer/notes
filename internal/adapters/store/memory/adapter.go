package memory

import (
	"context"
	"sort"
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
	now := timestamppb.Now()
	note := &notesv1.Note{
		Id:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
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
	return proto.Clone(note).(*notesv1.Note), nil
}

func (m *MemoryStore) ListNotes(_ context.Context, pageSize int32, offset int32) ([]*notesv1.Note, int32) {
	m.mu.RLock()
	// * Creating sorted list for consistency
	notesList := make([]*notesv1.Note, 0, len(m.data))
	for _, v := range m.data {
		notesList = append(notesList, proto.Clone(v).(*notesv1.Note))
	}
	sort.Slice(notesList, func(i, j int) bool {
		return notesList[i].CreatedAt.AsTime().After(notesList[j].CreatedAt.AsTime())
	})
	size := int32(len(notesList))
	if offset > size {
		offset = size
	}
	endPos := offset + pageSize // * Calculating last element for read
	if endPos >= size {
		return notesList[offset:], 0
	}
	return notesList[offset:endPos], endPos // * endPos - next elem position
}

func (m *MemoryStore) UpdateNote(_ context.Context, note *notesv1.Note) (*notesv1.Note, error) {
	m.mu.Lock() // * One mutex for atomicity
	_, ok := m.data[note.Id]
	if !ok {
		return nil, store.ErrStoreNotFound
	}
	m.data[note.Id] = proto.Clone(note).(*notesv1.Note)
	m.mu.Unlock()
	return proto.Clone(note).(*notesv1.Note), nil
}

func (m *MemoryStore) DeleteNote(_ context.Context, id string) {
	m.mu.Lock()
	delete(m.data, id)
	m.mu.Unlock()
}
