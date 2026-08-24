package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflow"
)

func main() {
	output := flag.String("out", "examples/rag-workflow/provider-free", "fixture output directory")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "fixture generator accepts no positional arguments")
		os.Exit(2)
	}
	_, err := ragworkflow.WriteProviderFreeFixtures(context.Background(), *output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fixture generation failed")
		os.Exit(1)
	}
	body, err := os.ReadFile(*output + "/manifest.json")
	if err != nil {
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(body); err != nil {
		os.Exit(1)
	}
}
