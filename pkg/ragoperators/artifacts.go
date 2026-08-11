package ragoperators

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func materializationData(records any) ([]byte, string) {
	data, _ := json.Marshal(records)
	sum := sha256.Sum256(data)
	return data, "sha256:" + hex.EncodeToString(sum[:])
}
func materializedArtifact(role, kind, name, schema string, data []byte, manifest any) Artifact {
	metadata, _ := json.Marshal(manifest)
	return Artifact{Role: role, Kind: kind, Name: name, SchemaVersion: schema, MediaType: "application/json", Metadata: metadata, Data: data}
}
func parent(role, digest, schema string) []ragcontract.ParentDigest {
	if digest == "" {
		return nil
	}
	return []ragcontract.ParentDigest{{Role: role, Digest: digest, SchemaVersion: schema}}
}
