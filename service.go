package todo

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/xid"
)

// TodoService for Todos
type TodoService interface {
	GetAllForUser(ctx context.Context, username string) ([]Todo, error)
	GetByID(ctx context.Context, id string) (Todo, error)
	Add(ctx context.Context, todo Todo) (Todo, error)
	Update(ctx context.Context, id string, todo Todo) error
	Delete(ctx context.Context, id string) error
}

// *** Implementation ***

var (
	// ErrInconsistentIDs is when the ID of the Entity you are updating differs from the ID given
	ErrInconsistentIDs = errors.New("Inconsistent IDs")
	// ErrNotFound is when the Entity doesn't exist
	ErrNotFound = errors.New("Not found")
)

// NewInmemTodoService creates an in memory Todo service
func NewInmemTodoService() TodoService {
	s := &inmemService{
		m: map[string]Todo{},
	}
	rand.Seed(time.Now().UnixNano())
	return s
}

// inmemService is a In Memory implementation of the service
type inmemService struct {
	sync.RWMutex
	m map[string]Todo
}

// GetAllForUser gets Todos from memory for a user
func (s *inmemService) GetAllForUser(ctx context.Context, username string) ([]Todo, error) {
	s.RLock()
	defer s.RUnlock()

	todos := make([]Todo, 0, len(s.m))
	for _, todo := range s.m {
		//if todo.Username == username {
		todos = append(todos, todo)
		//}
	}
	return todos, nil
}

// Get an Todos from the database
func (s *inmemService) GetByID(ctx context.Context, id string) (Todo, error) {
	s.Lock()
	defer s.Unlock()

	if todo, ok := s.m[id]; ok {
		return todo, nil
	}

	return Todo{}, ErrNotFound
}

// Add a Todo to memory
func (s *inmemService) Add(ctx context.Context, todo Todo) (Todo, error) {
	s.Lock()
	defer s.Unlock()

	todo.ID = xid.New().String()
	todo.CreatedOn = time.Now()

	s.m[todo.ID] = todo
	return todo, nil
}

// Update a Todo in memory
func (s *inmemService) Update(ctx context.Context, id string, todo Todo) error {
	s.Lock()
	defer s.Unlock()

	if id != todo.ID {
		return ErrInconsistentIDs
	}

	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}

	s.m[todo.ID] = todo
	return nil
}

// Delete a Todo from memory
func (s *inmemService) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}

	delete(s.m, id)
	return nil
}
