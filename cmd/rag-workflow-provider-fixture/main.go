package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflow"
)

func main() {
	output := flag.String("output", "examples/rag-workflow-provider", "output directory")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "fixture generator does not accept positional arguments")
		os.Exit(2)
	}
	manifest, err := ragworkflow.WriteDeterministicProviderFixtures(context.Background(), *output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "provider fixture generation failed")
		os.Exit(1)
	}
	fmt.Printf("generated %d provider fixture cases in %s\n", len(manifest.Cases), *output)
}
