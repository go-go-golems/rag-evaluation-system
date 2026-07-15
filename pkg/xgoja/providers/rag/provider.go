// Package rag registers the typed RAG laboratory module for xgoja/v2 hosts.
package rag

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	ragmodule "github.com/go-go-golems/rag-evaluation-system/pkg/gojamodules/rag"
)

const PackageID = "rag-evaluation-system"

// Register exposes require("rag") to generated xgoja binaries. The module has
// no provider configuration: each script explicitly opens a laboratory and
// chooses whether it is read-only or allowed to create experiment runs.
func Register(registry *providerapi.ProviderRegistry) error {
	return registry.Package(PackageID, providerapi.Module{
		Name:        ragmodule.ModuleName,
		DefaultAs:   ragmodule.ModuleName,
		Description: "Typed fluent RAG laboratory builders and explicit immutable experiment runs.",
		TypeScript:  ragmodule.TypeScriptModule(),
		NewModuleFactory: func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
			return ragmodule.NewLoader(), nil
		},
	})
}
