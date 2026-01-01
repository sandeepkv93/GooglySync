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
	LastEvent string
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

// Update replaces the current snapshot, preserving LastEvent when omitted.
func (s *Store) Update(snapshot Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now()
	}
	if snapshot.LastEvent == "" {
		snapshot.LastEvent = s.snapshot.LastEvent
	}
	s.snapshot = snapshot
}

// SetLastEvent updates the last observed event without changing state.
func (s *Store) SetLastEvent(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshot.LastEvent = msg
	s.snapshot.UpdatedAt = time.Now()
}

// Current returns a copy of the latest snapshot.
func (s *Store) Current() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.snapshot
}
