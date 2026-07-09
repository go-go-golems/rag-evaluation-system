package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl/migrationcheck"
)

func main() {
	jsonOutput := flag.Bool("json", false, "emit findings as JSON")
	failOnFindings := flag.Bool("fail-on-findings", false, "exit 1 when findings are reported")
	root := flag.String("root", "", "repository root used for relative output paths; defaults to current working directory")
	flag.Parse()

	findings, err := migrationcheck.ScanPaths(migrationcheck.Options{
		Root:  *root,
		Paths: flag.Args(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "widgetdsl-migration-checker: %v\n", err)
		os.Exit(2)
	}

	if *jsonOutput {
		data, err := migrationcheck.FindingsJSON(findings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "widgetdsl-migration-checker: encode JSON: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(string(data))
	} else if len(findings) == 0 {
		fmt.Println("No legacy Widget DSL imports or raw component escape hatches found.")
	} else {
		fmt.Printf("Found %d migration finding(s):\n", len(findings))
		for _, finding := range findings {
			fmt.Printf("%s:%d: %s %s: %s\n", finding.Path, finding.Line, finding.Kind, finding.Value, finding.Text)
		}
	}

	if *failOnFindings && len(findings) > 0 {
		os.Exit(1)
	}
}
