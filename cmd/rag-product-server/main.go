package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproduct"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/spf13/cobra"
)

type serverCommand struct{ *cmds.CommandDescription }

type serverSettings struct {
	Address         string `glazed:"address"`
	PlanPath        string `glazed:"plan"`
	CorpusPath      string `glazed:"corpus"`
	ProviderProfile string `glazed:"provider-profile"`
	ProviderConfig  string `glazed:"provider-config"`
}

var _ cmds.BareCommand = (*serverCommand)(nil)

func newServerCommand() (*cobra.Command, error) {
	description := cmds.NewCommandDescription(
		"rag-product-server",
		cmds.WithShort("Serve a canonical RAG product plan"),
		cmds.WithFlags(
			fields.New("address", fields.TypeString, fields.WithDefault("127.0.0.1:8780"), fields.WithHelp("Listen address")),
			fields.New("plan", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Canonical rag-product-plan/v2 JSON")),
			fields.New("corpus", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Verified rag-corpus-artifact JSON")),
			fields.New("provider-profile", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Explicit provider profile: fixtures or real")),
			fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only real-provider configuration YAML; required for provider-profile=real")),
		),
	)
	return cli.BuildCobraCommandFromCommand(&serverCommand{CommandDescription: description}, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-product-server", ShortHelpSections: []string{schema.DefaultSlug}}))
}

func (c *serverCommand) Run(ctx context.Context, parsed *values.Values) error {
	settings := &serverSettings{}
	if err := parsed.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	return runServer(ctx, settings)
}

func main() {
	command, err := newServerCommand()
	cobra.CheckErr(err)
	cobra.CheckErr(command.Execute())
}

func runServer(parent context.Context, settings *serverSettings) error {
	planFile, err := os.Open(settings.PlanPath)
	if err != nil {
		return err
	}
	defer func() { _ = planFile.Close() }()
	plan, err := ragproduct.Load(planFile)
	if err != nil {
		return err
	}
	corpusFile, err := os.Open(settings.CorpusPath)
	if err != nil {
		return err
	}
	defer func() { _ = corpusFile.Close() }()
	var corpus ragoperators.CorpusArtifact
	decoder := json.NewDecoder(corpusFile)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&corpus); err != nil {
		return err
	}
	bindings, closeProviders, err := productBindings(parent, settings, corpus)
	if err != nil {
		return err
	}
	defer closeProviders()
	runtime, err := ragproduct.New(parent, plan, bindings)
	if err != nil {
		return err
	}
	defer func() { _ = runtime.Close() }()
	server := &http.Server{Addr: settings.Address, Handler: newHandler(runtime), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Printf("rag product server listening on %s plan=%s", settings.Address, runtime.PlanID())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func productBindings(ctx context.Context, settings *serverSettings, corpus ragoperators.CorpusArtifact) (ragproduct.Bindings, func(), error) {
	switch settings.ProviderProfile {
	case "fixtures":
		if settings.ProviderConfig != "" {
			return ragproduct.Bindings{}, nil, fmt.Errorf("fixture profile does not accept provider config")
		}
		fixtures := ragoperators.NewFixtureProviders()
		return ragproduct.Bindings{Corpus: corpus, Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache()}, func() {}, nil
	case "real":
		if settings.ProviderConfig == "" {
			return ragproduct.Bindings{}, nil, fmt.Errorf("real profile requires provider config")
		}
		set, err := ragproviders.Load(ctx, settings.ProviderConfig)
		if err != nil {
			return ragproduct.Bindings{}, nil, err
		}
		return ragproduct.Bindings{Corpus: corpus, Manifests: set.Manifests, Schemas: set.Schemas, Generator: set.Generator, Embedder: set.Embedder, Reranker: set.Reranker, Cache: set.Cache}, func() { _ = set.Close() }, nil
	default:
		return ragproduct.Bindings{}, nil, fmt.Errorf("unsupported provider profile")
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
