package widgetdsl

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

// Registrar registers the single widget.dsl module with a go-go-goja engine runtime.
type Registrar struct{}

var _ engine.RuntimeModuleRegistrar = (*Registrar)(nil)

func NewRegistrar() *Registrar { return &Registrar{} }

func (r *Registrar) ID() string { return "widget-dsl" }

func (r *Registrar) RegisterRuntimeModule(_ *engine.RuntimeModuleRegistrationContext, reg *require.Registry) error {
	Register(reg)
	return nil
}
