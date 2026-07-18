package providers

import "github.com/spf13/cobra"

// NewCommand creates the Glazed-backed provider host command group.
func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "providers", Short: "Inspect and validate real-provider host configuration"}
	validate, err := NewValidateCommand()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(validate)
	return cmd, nil
}
