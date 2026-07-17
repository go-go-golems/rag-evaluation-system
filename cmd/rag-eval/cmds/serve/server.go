package serve

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/api"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/web"
	"github.com/rs/zerolog"
	zerolog_log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func runHTTPServer(cmd *cobra.Command, address, dbPath, engineDB, logLevel string) error {
	// Configure logging
	zerolog.SetGlobalLevel(parseLogLevel(logLevel))
	zerolog_log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	// Open database
	database, err := db.OpenDB(dbPath)
	if err != nil {
		zerolog_log.Fatal().Err(err).Str("path", dbPath).Msg("failed to open database")
	}
	defer func() { _ = database.Close() }()

	if err := db.Migrate(database); err != nil {
		zerolog_log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Wire HTTP handlers
	mux := http.NewServeMux()
	api.RegisterHandlersWithOptions(mux, database, api.Options{EngineDB: engineDB})
	mux.Handle("/", web.SPAHandler())

	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		zerolog_log.Info().Str("addr", address).Msg("starting RAG evaluation server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zerolog_log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zerolog_log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
