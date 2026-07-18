package study

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	ragmodule "github.com/go-go-golems/rag-evaluation-system/pkg/gojamodules/rag"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func LoadStudy(path string) (ragcontract.Study, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return ragcontract.Study{}, err
	}
	vm := goja.New()
	registry := require.NewRegistry()
	registry.RegisterNativeModule(ragmodule.ModuleName, ragmodule.NewLoader())
	registry.Enable(vm)
	encoded, _ := json.Marshal(filepath.ToSlash(absolute))
	value, err := vm.RunString("require(" + string(encoded) + ")")
	if err != nil {
		return ragcontract.Study{}, fmt.Errorf("RAG_STUDY_LOAD: %w", err)
	}
	raw, err := json.Marshal(value.Export())
	if err != nil {
		return ragcontract.Study{}, err
	}
	return ragcontract.DecodeStudy(bytes.NewReader(raw))
}
