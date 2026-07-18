// Package glazedcobra provides a transitional adapter for legacy Cobra command
// implementations. It makes Glazed the sole parser and command definition at
// the CLI boundary while command behavior is migrated incrementally.
package glazedcobra

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// WrapTree replaces every executable legacy Cobra leaf with a Glazed command.
// Groups remain structural Cobra nodes because Cobra owns hierarchical command
// dispatch; they do not define or parse flags.
func WrapTree(legacy *cobra.Command) (*cobra.Command, error) {
	if legacy == nil {
		return nil, fmt.Errorf("legacy command is nil")
	}
	if len(legacy.Commands()) == 0 {
		return wrapLeaf(legacy)
	}
	group := &cobra.Command{Use: legacy.Use, Aliases: legacy.Aliases, Short: legacy.Short, Long: legacy.Long, SilenceUsage: legacy.SilenceUsage, SilenceErrors: legacy.SilenceErrors}
	for _, child := range legacy.Commands() {
		wrapped, err := WrapTree(child)
		if err != nil {
			return nil, err
		}
		group.AddCommand(wrapped)
	}
	return group, nil
}

type positionalArgument struct {
	name     string
	required bool
}

type legacyCommand struct {
	*cmds.CommandDescription
	legacy    *cobra.Command
	arguments []positionalArgument
}

func wrapLeaf(legacy *cobra.Command) (*cobra.Command, error) {
	flags, err := glazedFlags(legacy.Flags())
	if err != nil {
		return nil, fmt.Errorf("%s flags: %w", legacy.Name(), err)
	}
	arguments := positionalArguments(legacy.Use)
	argumentFields := make([]*fields.Definition, 0, len(arguments))
	for _, argument := range arguments {
		argumentFields = append(argumentFields, fields.New(argument.name, fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(argument.required)))
	}
	description := cmds.NewCommandDescription(legacy.Name(), cmds.WithShort(legacy.Short), cmds.WithLong(legacy.Long), cmds.WithFlags(flags...), cmds.WithArguments(argumentFields...))
	command := &legacyCommand{CommandDescription: description, legacy: legacy, arguments: arguments}
	return cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{ShortHelpSections: []string{schema.DefaultSlug}, SkipCommandSettingsSection: true}))
}

func (c *legacyCommand) Run(ctx context.Context, parsed *values.Values) error {
	for name, value := range parsed.GetDataMap() {
		if c.legacy.Flags().Lookup(name) == nil {
			continue
		}
		if err := c.legacy.Flags().Set(name, pflagValue(value)); err != nil {
			return fmt.Errorf("set legacy flag %s: %w", name, err)
		}
	}
	args := make([]string, 0, len(c.arguments))
	for _, argument := range c.arguments {
		value, ok := parsed.GetDataMap()[argument.name]
		if !ok {
			if argument.required {
				return fmt.Errorf("missing argument %s", argument.name)
			}
			continue
		}
		args = append(args, pflagValue(value))
	}
	c.legacy.SetContext(ctx)
	c.legacy.SetOut(io.Discard)
	c.legacy.SetErr(io.Discard)
	if c.legacy.RunE != nil {
		return c.legacy.RunE(c.legacy, args)
	}
	if c.legacy.Run != nil {
		c.legacy.Run(c.legacy, args)
		return nil
	}
	return fmt.Errorf("legacy command %s has no runner", c.legacy.Name())
}

func glazedFlags(set *pflag.FlagSet) ([]*fields.Definition, error) {
	result := []*fields.Definition{}
	var conversionErr error
	set.VisitAll(func(flag *pflag.Flag) {
		if conversionErr != nil || flag.Hidden {
			return
		}
		typ, err := glazedType(flag.Value.Type())
		if err != nil {
			conversionErr = err
			return
		}
		options := []fields.Option{fields.WithDefault(glazedDefault(flag, typ)), fields.WithHelp(flag.Usage)}
		if flag.Shorthand != "" {
			options = append(options, fields.WithShortFlag(flag.Shorthand))
		}
		result = append(result, fields.New(flag.Name, typ, options...))
	})
	return result, conversionErr
}

func glazedDefault(flag *pflag.Flag, typ fields.Type) any {
	switch typ {
	case fields.TypeBool:
		value, _ := strconv.ParseBool(flag.DefValue)
		return value
	case fields.TypeInteger:
		value, _ := strconv.ParseInt(flag.DefValue, 10, 64)
		return value
	case fields.TypeFloat:
		value, _ := strconv.ParseFloat(flag.DefValue, 64)
		return value
	default:
		return flag.DefValue
	}
}

func glazedType(pflagType string) (fields.Type, error) {
	switch pflagType {
	case "string", "duration":
		return fields.TypeString, nil
	case "bool":
		return fields.TypeBool, nil
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return fields.TypeInteger, nil
	case "float32", "float64":
		return fields.TypeFloat, nil
	case "stringSlice", "stringArray":
		// Keep pflag's native CSV/JSON parsing semantics in the delegated handler.
		return fields.TypeString, nil
	default:
		return "", fmt.Errorf("unsupported pflag type %q", pflagType)
	}
}

func positionalArguments(use string) []positionalArgument {
	result := []positionalArgument{}
	for _, word := range strings.Fields(use) {
		switch {
		case strings.HasPrefix(word, "<") && strings.HasSuffix(word, ">"):
			result = append(result, positionalArgument{name: strings.TrimSuffix(strings.TrimPrefix(word, "<"), ">"), required: true})
		case strings.HasPrefix(word, "[") && strings.HasSuffix(word, "]"):
			result = append(result, positionalArgument{name: strings.TrimSuffix(strings.TrimPrefix(word, "["), "]")})
		}
	}
	return result
}

func pflagValue(value any) string {
	switch typed := value.(type) {
	case []string:
		return strings.Join(typed, ",")
	default:
		return fmt.Sprint(value)
	}
}
