package server

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Morditux/serverlib/templates"
)

// ServerInstance represents the singleton instance of the server.
var ServerInstance *Server

// Server represents the HTTP server.
type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	t          *templates.Templates
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
			Address: ":8080",
			Handler: mux,
		}
	} else {
		serverConfig = config[0]
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

		router: serverConfig.Handler.(*http.ServeMux),
	}
	return ServerInstance
}

// Start starts the server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop stops the server.
func (s *Server) Stop() error {
	return s.httpServer.Close()
}

// HandleFunc registers a function to handle HTTP requests with the given pattern.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// Handle registers a handler to handle HTTP requests with the given pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, handler)
}

// AddTemplateSource adds a new template source to the server's template manager.
// The source parameter specifies the template source path to be added.
func (s *Server) AddTemplateSource(source string) {
	s.t.AddSource(source)
}

// Render renders the specified template with the given data and writes the result to the response writer.
func (s *Server) Render(w io.Writer, template string, data map[string]interface{}) {
	s.t.Execute(w, template, data)
}

func (s *Server) Templates() *templates.Templates {
	return s.t
}
