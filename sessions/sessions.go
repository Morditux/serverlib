package sessions

// Session represents a user session with methods to manage session data.
// It provides an interface for retrieving and storing key-value pairs.
//
// Methods:
//   - Id() string: Returns the unique identifier for the session.
//   - Get(key string) (any, bool): Retrieves the value associated with the given key.
//     Returns the value and a boolean indicating whether the key was found.
//   - Set(key string, value any): Stores the given value associated with the key.
type Session interface {
	// Id returns the unique identifier for the session.
	Id() string
	// Get retrieves the value associated with the given key.
	// Returns the value or nil if the key was not found.
	Get(key string) any
	// Set stores the given value associated with the key.
	Set(key string, value any)
	// Check if the key exists in the session.
	Exists(key string) bool
}

// Sessions defines an interface for managing user sessions.
// It provides methods to retrieve, store, and delete sessions by their unique identifier.
//
// Methods:
//   - Get(id string) (Session, bool): Retrieves a session by its ID. Returns the session is found.
//   - Set(id string, session Session): Stores a session with the given ID.
//   - Delete(id string): Deletes the session associated with the given ID.
type Sessions interface {
	// Get retrieves a session by its ID.
	// Returns the session and a boolean indicating whether the session was found.
	Get(id string) (Session, bool)
	// Set stores a session with the given ID.
	Set(id string, session Session)
	// Delete deletes the session associated with the given ID.
	Delete(id string)
	// Create a new session with a new ID.
	New() Session
}
