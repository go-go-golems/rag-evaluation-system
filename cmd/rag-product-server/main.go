package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproduct"
)

func main() {
	var address, planPath, corpusPath string
	flag.StringVar(&address, "address", "127.0.0.1:8780", "listen address")
	flag.StringVar(&planPath, "plan", "", "canonical rag-product-plan/v2 JSON")
	flag.StringVar(&corpusPath, "corpus", "", "verified rag-corpus-artifact JSON")
	flag.Parse()
	if planPath == "" || corpusPath == "" {
		log.Fatal("--plan and --corpus are required")
	}
	planFile, err := os.Open(planPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = planFile.Close() }()
	plan, err := ragproduct.Load(planFile)
	if err != nil {
		log.Fatal(err)
	}
	corpusFile, err := os.Open(corpusPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = corpusFile.Close() }()
	var corpus ragoperators.CorpusArtifact
	decoder := json.NewDecoder(corpusFile)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&corpus); err != nil {
		log.Fatal(err)
	}
	fixtures := ragoperators.NewFixtureProviders()
	runtime, err := ragproduct.New(context.Background(), plan, ragproduct.Bindings{Corpus: corpus, Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Addr: address, Handler: newHandler(runtime), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Printf("rag product server listening on %s plan=%s", address, runtime.PlanID())
	err = server.ListenAndServe()
	_ = runtime.Close()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newHandler(runtime *ragproduct.Runtime) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "planId": runtime.PlanID()})
	})
	mux.HandleFunc("POST /v1/query", func(w http.ResponseWriter, request *http.Request) {
		if contentType := request.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
			writeError(w, http.StatusUnsupportedMediaType, "RAG_PRODUCT_CONTENT_TYPE")
			return
		}
		body := http.MaxBytesReader(w, request.Body, 1<<20)
		productRequest, err := ragproduct.DecodeRequest(body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "RAG_PRODUCT_REQUEST_INVALID")
			return
		}
		response, err := runtime.Execute(request.Context(), productRequest)
		if err != nil {
			status := http.StatusInternalServerError
			if strings.Contains(err.Error(), "REQUEST_") {
				status = http.StatusBadRequest
			}
			if request.Context().Err() != nil {
				status = http.StatusRequestTimeout
			}
			writeError(w, status, "RAG_PRODUCT_EXECUTION_FAILED")
			return
		}
		writeJSON(w, http.StatusOK, response)
	})
	return mux
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"code": code})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "encoding failure", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(data, '\n'))
}
