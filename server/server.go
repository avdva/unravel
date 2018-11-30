// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/avdva/unravel/card"
)

// Server is a HTTP server, that handles JSON requests.
type Server struct {
	srv     *http.Server
	handler card.Handler
}

// New returns new Server.
func New(addr string, handler card.Handler) *Server {
	result := &Server{
		handler: handler,
	}
	mux := http.NewServeMux()
	mux.Handle("/api/card", chain(result, methodMiddleware, corsMiddleware))
	result.srv = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return result
}

// Serve starts processing loop.
func (s *Server) Serve() error {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop shutdowns HTTP server.
func (s *Server) Stop() {
	s.srv.Shutdown(context.Background())
}

// ServeHTTP is a HTTP requests handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := decodeRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: " + err.Error()))
		return
	}
	s.handleRequest(req)
}

func (s *Server) handleRequest(req *request) {
	header := card.EventHeader{
		WebsiteUrl: req.URL,
		SessionID:  req.SessionID,
	}
	switch req.EventType {
	case evCopyPaste:
		if req.FormId != nil && req.Pasted != nil {
			s.handler.OnCopyPaste(header, *req.FormId, *req.Pasted)
		}
	case evWindowResize:
		if req.ResizeFrom != nil && req.ResizeTo != nil {
			from := card.Dimension(*req.ResizeFrom)
			to := card.Dimension(*req.ResizeTo)
			s.handler.OnResize(header, from, to)
		}
	case evSubmit:
		if req.Time != nil {
			s.handler.OnSubmit(header, *req.Time)
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		next.ServeHTTP(w, r)
	})
}

func methodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			next.ServeHTTP(w, r)
		}
	})
}

func chain(h http.Handler, mids ...func(next http.Handler) http.Handler) http.Handler {
	for _, m := range mids {
		h = m(h)
	}
	return h
}
