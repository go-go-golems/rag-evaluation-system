// Builds the minimal canonical researchctl specification used only to verify
// the fixture sweep's compact operation-custody export/import path.
package main

import (
	"fmt"

	"github.com/go-go-golems/researchctl/pkg/lab"
)

func main() {
	identity := lab.ExecutionIdentity{
		SchemaVersion:       lab.ExecutionSpecSchemaVersion,
		IdentityScheme:      lab.ExecutionIdentityScheme,
		Domain:              "rag",
		DomainSchemaVersion: "rag/v2",
		DomainConfig:        []byte(`{"schemaVersion":"rag-ttc-v3-sweep-custody/v1"}`),
	}
	id, _, err := lab.ExecutionID(identity)
	if err != nil {
		panic(err)
	}
	body, err := lab.CanonicalJSON(lab.SpecificationRecord{ID: id, IdentityScheme: lab.ExecutionIdentityScheme, CanonicalIdentity: identity, DisplayName: "fixture TTC custody"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
