package widgetserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/defaultspa"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetrunner"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetschema"
)

type FrontendMode string

const (
	FrontendEmbedded FrontendMode = "embedded"
	FrontendDir      FrontendMode = "dir"
	FrontendProxy    FrontendMode = "proxy"
	FrontendAPIOnly  FrontendMode = "api-only"
)

type Config struct {
	Addr            string
	Runner          *widgetrunner.Runner
	APIPrefix       string
	Dev             bool
	FrontendMode    FrontendMode
	FrontendDir     string
	FrontendProxy   string
	FrontendHandler http.Handler
}

type Server struct {
	cfg Config
	mux *http.ServeMux
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func New(cfg Config) (*Server, error) {
	if cfg.Runner == nil {
		return nil, fmt.Errorf("widgetserver: runner is required")
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.APIPrefix == "" {
		cfg.APIPrefix = "/api/widget"
	}
	cfg.APIPrefix = cleanPrefix(cfg.APIPrefix)
	if cfg.FrontendMode == "" {
		cfg.FrontendMode = FrontendEmbedded
	}

	s := &Server{cfg: cfg, mux: http.NewServeMux()}
	if err := s.registerRoutes(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) Run(ctx context.Context) error {
	httpSrv := &http.Server{Addr: s.cfg.Addr, Handler: s.Handler(), ReadHeaderTimeout: 5 * time.Second}
	errCh := make(chan error, 1)
	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) registerRoutes() error {
	prefix := s.cfg.APIPrefix
	s.mux.HandleFunc("GET "+prefix+"/health", s.handleHealth)
	s.mux.HandleFunc("GET "+prefix+"/pages/{id}", s.handlePage)
	s.mux.HandleFunc("POST "+prefix+"/actions/{name}", s.handleAction)
	s.mux.HandleFunc("GET "+prefix+"/schema", s.handleSchema)

	if s.cfg.FrontendMode != FrontendAPIOnly {
		handler, err := s.frontendHandler()
		if err != nil {
			return err
		}
		s.mux.Handle("/", handler)
	}
	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "widget-site"})
}

func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	page, err := s.cfg.Runner.RenderPage(r.Context(), id, widgetrunner.PageContext{Query: queryMap(r.URL.Query())})
	if err != nil {
		if errors.Is(err, widgetrunner.ErrPageNotFound) {
			writeError(w, http.StatusNotFound, "page_not_found", err.Error(), map[string]any{"pageId": id})
			return
		}
		if errors.Is(err, widgetrunner.ErrInvalidPage) || errors.Is(err, widgetrunner.ErrInvalidWidgetIR) {
			writeError(w, http.StatusBadRequest, "invalid_widget_ir", err.Error(), map[string]any{"pageId": id})
			return
		}
		writeError(w, http.StatusInternalServerError, "script_runtime_error", scriptErrorMessage(err, s.cfg.Dev), map[string]any{"pageId": id})
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var req widgetrunner.ActionRequest
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "invalid_action_request", err.Error(), map[string]any{"action": name})
			return
		}
	}
	result, err := s.cfg.Runner.InvokeAction(r.Context(), name, req)
	if err != nil {
		if errors.Is(err, widgetrunner.ErrActionNotFound) {
			writeError(w, http.StatusNotFound, "action_not_found", err.Error(), map[string]any{"action": name})
			return
		}
		if errors.Is(err, widgetrunner.ErrInvalidAction) {
			writeError(w, http.StatusBadRequest, "invalid_action_result", err.Error(), map[string]any{"action": name})
			return
		}
		writeError(w, http.StatusInternalServerError, "action_runtime_error", scriptErrorMessage(err, s.cfg.Dev), map[string]any{"action": name})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleSchema(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, widgetschema.Summary())
}

func (s *Server) frontendHandler() (http.Handler, error) {
	switch s.cfg.FrontendMode {
	case FrontendAPIOnly:
		return http.NotFoundHandler(), nil
	case FrontendEmbedded:
		if s.cfg.FrontendHandler != nil {
			return s.cfg.FrontendHandler, nil
		}
		return defaultspa.Handler(), nil
	case FrontendDir:
		if strings.TrimSpace(s.cfg.FrontendDir) == "" {
			return nil, fmt.Errorf("widgetserver: dir frontend mode requires FrontendDir")
		}
		return SPAHandler(os.DirFS(s.cfg.FrontendDir)), nil
	case FrontendProxy:
		if strings.TrimSpace(s.cfg.FrontendProxy) == "" {
			return nil, fmt.Errorf("widgetserver: proxy frontend mode requires FrontendProxy")
		}
		target, err := url.Parse(s.cfg.FrontendProxy)
		if err != nil {
			return nil, fmt.Errorf("widgetserver: parse frontend proxy URL: %w", err)
		}
		return httputil.NewSingleHostReverseProxy(target), nil
	default:
		return nil, fmt.Errorf("widgetserver: unsupported frontend mode %q", s.cfg.FrontendMode)
	}
}

func SPAHandler(publicFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(publicFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if f, err := publicFS.Open(path); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		index, err := fs.ReadFile(publicFS, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(index)
	})
}

func queryMap(values url.Values) map[string]string {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]string, len(values))
	for key, value := range values {
		if len(value) > 0 {
			out[key] = value[0]
		} else {
			out[key] = ""
		}
	}
	return out
}

func cleanPrefix(prefix string) string {
	prefix = "/" + strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "/" {
		return "/api/widget"
	}
	return strings.TrimRight(prefix, "/")
}

func scriptErrorMessage(err error, dev bool) string {
	if dev {
		return err.Error()
	}
	return "internal server error"
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	writeJSON(w, status, ErrorResponse{Error: ErrorBody{Code: code, Message: message, Details: details}})
}
