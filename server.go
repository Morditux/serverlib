package serverlib

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Morditux/serverlib/server/sessions"
	"github.com/Morditux/serverlib/templates"
	"github.com/google/uuid"
)

type LogLevel int

const (
	None LogLevel = iota
	Info
	Debug
	Error
)

// ServerInstance represents the singleton instance of the server.
var ServerInstance *Server

// Server represents an HTTP server with routing and session management capabilities.
// It includes an HTTP server, a router for handling HTTP requests, a session manager
// for managing user sessions, a session key for session security, and a template
// engine for rendering HTML templates.
type Server struct {
	httpServer     *http.Server
	router         *http.ServeMux
	sessionManager sessions.Sessions
	sessionKey     string
	logger         *log.Logger
	dateFormat     func(time.Time) string
	t              *templates.Templates
	logLevel       LogLevel
}

type ServerConfig struct {
	Address                      string
	Handler                      http.Handler
	DisableGeneralOptionsHandler bool
	TLSConfig                    *tls.Config
	ReadTimeout                  time.Duration
	ReadHeaderTimeout            time.Duration
	WriteTimeout                 time.Duration
	IdleTimeout                  time.Duration
	MaxHeaderBytes               int
	ConnState                    func(net.Conn, http.ConnState)
	ErrorLog                     *log.Logger
	BaseContext                  func(net.Listener) context.Context
	ConnContext                  func(ctx context.Context, c net.Conn) context.Context
	SessionManager               sessions.Sessions
	SessionKey                   string
	DateFormat                   func(time.Time) string
	LogLevel                     LogLevel
}

// NewServer creates a new instance of Server with the provided configuration.
// If no configuration is provided, it uses default settings with an address of ":8080" and a new ServeMux as the handler.
//
// Parameters:
//   - config: Optional variadic parameter of type ServerConfig. If provided, the first element is used as the server configuration.
//
// Returns:
//   - *Server: A pointer to the newly created Server instance.
func NewServer(config ...ServerConfig) *Server {
	var serverConfig ServerConfig

	if len(config) == 0 {
		mux := http.NewServeMux()
		serverConfig = ServerConfig{
			Address:        ":8080",
			Handler:        mux,
			SessionManager: sessions.NewMemorySessions(),
			SessionKey:     uuid.New().String(),
		}
	} else {
		serverConfig = config[0]
	}
	if serverConfig.SessionKey == "" {
		serverConfig.SessionKey = uuid.New().String()
	}
	if serverConfig.SessionManager == nil {
		serverConfig.SessionManager = sessions.NewMemorySessions()
	}
	if serverConfig.Handler == nil {
		serverConfig.Handler = http.NewServeMux()
	}
	if serverConfig.Address == "" {
		serverConfig.Address = ":8080"
	}
	if serverConfig.ErrorLog == nil {
		serverConfig.ErrorLog = log.New(os.Stderr, "", log.LstdFlags)
	}
	if serverConfig.DateFormat == nil {
		serverConfig.DateFormat = func(t time.Time) string {
			return t.Format(time.ANSIC)
		}
	}
	ServerInstance = &Server{
		t: templates.NewTemplates(),
		httpServer: &http.Server{
			Addr:              serverConfig.Address,
			Handler:           serverConfig.Handler,
			TLSConfig:         serverConfig.TLSConfig,
			ReadTimeout:       serverConfig.ReadTimeout,
			ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
			WriteTimeout:      serverConfig.WriteTimeout,
			IdleTimeout:       serverConfig.IdleTimeout,
			MaxHeaderBytes:    serverConfig.MaxHeaderBytes,
			ConnState:         serverConfig.ConnState,
			ErrorLog:          serverConfig.ErrorLog,
			BaseContext:       serverConfig.BaseContext,
			ConnContext:       serverConfig.ConnContext,
		},
		router:         serverConfig.Handler.(*http.ServeMux),
		sessionManager: serverConfig.SessionManager,
		sessionKey:     serverConfig.SessionKey,
		logger:         serverConfig.ErrorLog,
		dateFormat:     serverConfig.DateFormat,
		logLevel:       serverConfig.LogLevel,
	}
	return ServerInstance
}

// Start starts the server.
func (s *Server) Start() error {
	slog.Info("Server started", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop stops the server.
func (s *Server) Stop() error {
	slog.Info("Server stopped", "address", s.httpServer.Addr)
	return s.httpServer.Close()
}

// HandleFunc registers a function to handle HTTP requests with the given pattern.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	slog.Info("Registred HandleFunc", "pattern", pattern)
	s.router.HandleFunc(pattern, handler)
}

// Handle registers a handler to handle HTTP requests with the given pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	slog.Info("Registred handle", "pattern", pattern)
	s.router.Handle(pattern, handler)
}

// AddTemplateSource adds a new template source to the server's template manager.
// The source parameter specifies the template source path to be added.
func (s *Server) AddTemplateSource(source string) {
	slog.Info("Adding template source", "source", source)
	s.t.AddSource(source)
}

// Render renders the specified template with the given data and writes the result to the response writer.
func (s *Server) Render(w io.Writer, template string, data map[string]interface{}) {
	slog.Info("Rendering template", "template", template)
	s.t.Execute(w, template, data)
}

// Templates returns the server's templates.
// It provides access to the templates associated with the server instance.
func (s *Server) Templates() *templates.Templates {
	return s.t
}

// SessionManager returns the server's session manager.
// It provides access to the session manager associated with the server instance.
func (s *Server) Sessions() sessions.Sessions {
	return s.sessionManager
}

// SessionKey returns the session key associated with the server instance.
// This key is used to identify and manage user sessions.
func (s *Server) SessionKey() string {
	return s.sessionKey
}

// GetSession retrieves the session associated with the request's cookie.
// If the session does not exist, a new session is created and a new cookie is set.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//
// Returns:
//   - sessions.Session: The session associated with the request.
//   - bool: A boolean indicating whether the session was retrieved (true) or newly created (false).
func (s *Server) GetSession(w http.ResponseWriter, r *http.Request) (sessions.Session, bool) {
	cookie, err := r.Cookie(s.sessionKey)
	if err != nil {
		return nil, false
	}
	sessionID := cookie.Value
	session, ok := s.sessionManager.Get(sessionID)
	if !ok {
		// Create a new session if the session ID is not found
		sessionID = uuid.New().String()
		session = sessions.NewMemorySession(sessionID)
		s.sessionManager.Set(sessionID, session)
		http.SetCookie(w, &http.Cookie{
			Name:     s.sessionKey,
			Value:    sessionID,
			HttpOnly: true,
			MaxAge:   3600 * 24 * 7, // 1 week
		})
	}
	return session, ok
}

// SetLogLevel sets the logging level for the server.
//
// Parameters:
//
//	level (LogLevel): The desired logging level.
//
// Usage:
//
//	server.SetLogLevel(LogLevelDebug)
func (s *Server) SetLogLevel(level LogLevel) {
	s.logLevel = level
}

// LogInfo logs an informational message if the server's log level is set to Info or higher.
// It takes two parameters:
// - message: A string representing the message to be logged.
// - value: A string representing additional information to be logged alongside the message.
func LogInfo(message string, value string) {
	if ServerInstance.logLevel >= Info {
		ServerInstance.logger.Printf("INFO - %s: %s\n", message, value)
	}
}

// LogDebug logs a debug message if the server's log level is set to Debug or higher.
// It takes two parameters:
// - message: A string representing the debug message.
// - value: A string representing additional information to log with the message.
func LogDebug(message string, value string) {
	if ServerInstance.logLevel >= Debug {
		ServerInstance.logger.Printf("DEBUG - %s: %s\n", message, value)
	}
}

// LogError logs an error message with a specified value if the server's log level is set to Error or higher.
// Parameters:
//   - message: A string representing the error message to be logged.
//   - value: A string representing additional information or context about the error.
func LogError(message string, value string) {
	if ServerInstance.logLevel >= Error {
		ServerInstance.logger.Printf("ERROR - %s: %s\n", message, value)
	}
}
