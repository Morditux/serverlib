package sessions

import "sync"

// MemorySession represents an in-memory session with a unique identifier,
// a map to store session data, and a read-write mutex for concurrent access control.
type MemorySession struct {
	id   string
	data map[string]any
	mut  *sync.RWMutex
}

// MemorySessions is a struct that manages a collection of in-memory sessions.
// It contains a map of session IDs to MemorySession pointers and a read-write mutex
// to ensure thread-safe access to the sessions map.
type MemorySessions struct {
	sessions map[string]*MemorySession
	mut      *sync.RWMutex
}

// NewMemorySessions creates and returns a new instance of MemorySessions.
// It initializes the sessions map and the read-write mutex.
func NewMemorySessions() *MemorySessions {
	return &MemorySessions{
		sessions: make(map[string]*MemorySession),
		mut:      &sync.RWMutex{},
	}
}

// Get retrieves a session from the memory store by its ID.
// It returns the session and a boolean indicating whether the session was found.
// The method is thread-safe, using a read lock to ensure concurrent access.
func (s *MemorySessions) Get(id string) (Session, bool) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	session, ok := s.sessions[id]
	return session, ok
}

// Set stores a session in the MemorySessions map with the given id.
// It locks the mutex to ensure thread safety before modifying the map.
//
// Parameters:
//   - id: A string representing the session ID.
//   - session: A Session interface that will be type asserted to *MemorySession.
func (s *MemorySessions) Set(id string, session Session) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.sessions[id] = session.(*MemorySession)
}

// Delete removes a session from the memory store by its ID.
// It locks the session map to ensure thread safety during the deletion process.
//
// Parameters:
//
//	id (string): The ID of the session to be deleted.
func (s *MemorySessions) Delete(id string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.sessions, id)
}

// NewMemorySession creates a new MemorySession with the given id.
// It initializes the session data as an empty map and sets up a read-write mutex for concurrent access.
//
// Parameters:
//   - id: A string representing the unique identifier for the session.
//
// Returns:
//   - A pointer to a newly created MemorySession instance.
func NewMemorySession(id string) *MemorySession {
	return &MemorySession{
		id:   id,
		data: make(map[string]any),
		mut:  &sync.RWMutex{},
	}
}

// Id returns the unique identifier of the MemorySession.
// It retrieves the session's ID as a string.
func (s *MemorySession) Id() string {
	return s.id
}

// Get retrieves the value associated with the given key from the memory session.
// It returns the value and a boolean indicating whether the key was found.
// The method is thread-safe, using a read lock to ensure concurrent access.
func (s *MemorySession) Get(key string) (any, bool) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Set stores a key-value pair in the memory session. It locks the session
// to ensure thread safety before setting the value and unlocks it afterward.
//
// Parameters:
//   - key: The key under which the value will be stored.
//   - value: The value to be stored, which can be of any type.
func (s *MemorySession) Set(key string, value any) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.data[key] = value
}
