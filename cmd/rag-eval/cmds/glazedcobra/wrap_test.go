package glazedcobra

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestWrapTreeUsesGlazedParserAndDelegatesBehavior(t *testing.T) {
	var name string
	var count int
	var received []string
	legacy := &cobra.Command{Use: "demo <document> [label]", RunE: func(_ *cobra.Command, args []string) error {
		received = args
		if name != "Ada" || count != 3 {
			t.Fatalf("legacy values = %q, %d", name, count)
		}
		return nil
	}}
	legacy.Flags().StringVar(&name, "name", "", "name")
	legacy.Flags().IntVar(&count, "count", 0, "count")
	wrapped, err := WrapTree(legacy)
	if err != nil {
		t.Fatal(err)
	}
	wrapped.SetArgs([]string{"--name", "Ada", "--count", "3", "report"})
	if err := wrapped.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(received) != 1 || received[0] != "report" {
		t.Fatalf("args = %#v", received)
	}
}
