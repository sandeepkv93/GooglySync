package status

import (
	"sync"
	"time"
)

// State describes high-level sync state.
type State int

const (
	StateUnspecified State = iota
	StateIdle
	StateSyncing
	StateError
	StatePaused
)

// Snapshot captures current status.
type Snapshot struct {
	State     State
	Message   string
	UpdatedAt time.Time
}

// Store holds the latest status snapshot.
type Store struct {
	mu       sync.Mutex
	snapshot Snapshot
}

// NewStore constructs a status store with an initial idle state.
func NewStore() *Store {
	return &Store{
		snapshot: Snapshot{
			State:     StateIdle,
			Message:   "idle",
			UpdatedAt: time.Now(),
		},
	}
}

// Update replaces the current snapshot.
func (s *Store) Update(snapshot Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now()
	}
	s.snapshot = snapshot
}

// Current returns a copy of the latest snapshot.
func (s *Store) Current() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.snapshot
}
