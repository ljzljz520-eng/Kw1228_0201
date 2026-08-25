package httpapi

import (
	"context"
	"net/http"

	"example.com/energycore/internal/flow017"
)

type Server struct {
	runtime *flow017.Runtime
	mux     *http.ServeMux
}

func New(runtime *flow017.Runtime) *Server {
	server := &Server{runtime: runtime, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler {
	return requestID(s.logging(s.mux))
}

func (s *Server) Serve(ctx context.Context, address string) error {
	listener := &http.Server{Addr: address, Handler: s.Handler()}
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()
	err := listener.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/cores", s.cores)
	s.mux.HandleFunc("/imports", s.imports)
	s.mux.HandleFunc("/reports/summary", s.summary)
	s.adminRoutes()
}
