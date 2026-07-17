package search

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	cmds2 "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type SmokeCommand struct{ *cmds.CommandDescription }

var _ cmds.GlazeCommand = (*SmokeCommand)(nil)

type SmokeSettings struct {
	DB           string `glazed:"db"`
	IndexRoot    string `glazed:"index-root"`
	IndexID      string `glazed:"index-id"`
	File         string `glazed:"file"`
	Retriever    string `glazed:"retriever"`
	Limit        int    `glazed:"limit"`
	PreviewRunes int    `glazed:"preview-runes"`
}

type smokeFile struct {
	Queries []smokeQuery `yaml:"queries"`
}

type smokeQuery struct {
	ID                string   `yaml:"id"`
	Text              string   `yaml:"text"`
	Intent            string   `yaml:"intent"`
	ExpectedTerms     []string `yaml:"expected_terms"`
	ExpectedSourceIDs []string `yaml:"expected_source_ids"`
	Notes             string   `yaml:"notes"`
}

func newSmokeCommand() *cobra.Command {
	glazedCmd, err := newSmokeGlazeCommand()
	cobra.CheckErr(err)
	cobraCmd, err := cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCmd
}

func newSmokeGlazeCommand() (*SmokeCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &SmokeCommand{CommandDescription: cmds.NewCommandDescription(
		"smoke",
		cmds.WithShort("Run lightweight retrieval smoke checks from YAML"),
		cmds.WithLong(`Run a small manual query file against BM25 retrieval. This is a plumbing smoke check, not a benchmark.`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("index-root", fields.TypeString, fields.WithDefault(searchservice.DefaultIndexRoot), fields.WithHelp("Root directory for BM25 indexes")),
			fields.New("index-id", fields.TypeString, fields.WithHelp("BM25 index ID to query")),
			fields.New("file", fields.TypeString, fields.WithDefault("eval/ttc-smoke.yaml"), fields.WithHelp("Smoke query YAML file")),
			fields.New("retriever", fields.TypeString, fields.WithDefault("bm25"), fields.WithHelp("Retriever to smoke test; currently bm25")),
			fields.New("limit", fields.TypeInteger, fields.WithDefault(10), fields.WithHelp("Top-K retrieval limit per query")),
			fields.New("preview-runes", fields.TypeInteger, fields.WithDefault(160), fields.WithHelp("Preview runes per result")),
		), cmds.WithSections(glazedSection),
	)}, nil
}

func (c *SmokeCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	s := &SmokeSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	if s.Retriever != "bm25" {
		return fmt.Errorf("unsupported smoke retriever %q: currently only bm25 is implemented", s.Retriever)
	}
	contents, err := os.ReadFile(s.File)
	if err != nil {
		return err
	}
	var suite smokeFile
	if err := yaml.Unmarshal(contents, &suite); err != nil {
		return err
	}
	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()
	service := searchservice.NewService(queries, s.IndexRoot)
	for _, q := range suite.Queries {
		result, err := service.QueryBM25(ctx, searchservice.QueryRequest{IndexID: s.IndexID, Query: q.Text, Limit: s.Limit, PreviewRunes: s.PreviewRunes})
		status := "pass"
		message := ""
		if err != nil {
			status, message = "fail", err.Error()
		}
		matchedTerms := []string{}
		matchedSource := false
		if err == nil {
			if len(result.Items) == 0 {
				status, message = "fail", "no results"
			}
			for _, item := range result.Items {
				haystack := strings.ToLower(item.Title + " " + item.Preview)
				for _, term := range q.ExpectedTerms {
					termLower := strings.ToLower(term)
					if strings.Contains(haystack, termLower) && !contains(matchedTerms, termLower) {
						matchedTerms = append(matchedTerms, termLower)
					}
				}
				if len(q.ExpectedSourceIDs) == 0 || contains(q.ExpectedSourceIDs, item.SourceID) {
					matchedSource = true
				}
			}
			if status == "pass" && len(q.ExpectedTerms) > 0 && len(matchedTerms) == 0 {
				status, message = "warn", "no expected terms found in top results"
			}
			if status == "pass" && len(q.ExpectedSourceIDs) > 0 && !matchedSource {
				status, message = "warn", "no expected source found in top results"
			}
		}
		topChunk, topTitle := "", ""
		if err == nil && len(result.Items) > 0 {
			topChunk, topTitle = result.Items[0].ChunkID, result.Items[0].Title
		}
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("id", q.ID), types.MRP("query", q.Text), types.MRP("intent", q.Intent), types.MRP("retriever", s.Retriever), types.MRP("status", status), types.MRP("message", message), types.MRP("result_count", lenOrZero(result)), types.MRP("matched_terms", strings.Join(matchedTerms, ",")), types.MRP("top_chunk_id", topChunk), types.MRP("top_title", topTitle),
		)); err != nil {
			return err
		}
	}
	return nil
}

func contains(items []string, needle string) bool {
	for _, item := range items {
		if item == needle {
			return true
		}
	}
	return false
}

func lenOrZero(result *searchservice.QueryResult) int {
	if result == nil {
		return 0
	}
	return len(result.Items)
}
