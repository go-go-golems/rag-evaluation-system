package ragcontract

import (
	"fmt"
	"regexp"
)

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

func ValidateManifestBase(base ManifestBase, expectedSchema string, requireParent bool) error {
	if base.SchemaVersion != expectedSchema {
		return fmt.Errorf("RAG_V2_MANIFEST_SCHEMA: got %q, expected %q", base.SchemaVersion, expectedSchema)
	}
	if !digestPattern.MatchString(base.Digest) {
		return fmt.Errorf("RAG_V2_MANIFEST_DIGEST: invalid digest %q", base.Digest)
	}
	if requireParent && len(base.Parents) == 0 {
		return fmt.Errorf("RAG_V2_MANIFEST_PARENT: %s requires at least one parent digest", expectedSchema)
	}
	roles := map[string]bool{}
	for i, parent := range base.Parents {
		if parent.Role == "" || parent.SchemaVersion == "" || !digestPattern.MatchString(parent.Digest) {
			return fmt.Errorf("RAG_V2_MANIFEST_PARENT: invalid parent %d", i)
		}
		if roles[parent.Role] {
			return fmt.Errorf("RAG_V2_MANIFEST_PARENT: duplicate role %q", parent.Role)
		}
		roles[parent.Role] = true
	}
	if requireParent && base.Production == nil {
		return fmt.Errorf("RAG_V2_MANIFEST_PRODUCTION: %s requires producing operator/config", expectedSchema)
	}
	if base.Production != nil {
		if !ValidIdentifier(base.Production.Operator.Kind) || base.Production.Operator.Version == "" {
			return fmt.Errorf("RAG_V2_MANIFEST_PRODUCTION: invalid operator")
		}
		if _, err := CanonicalRaw(base.Production.Config, "{}"); err != nil {
			return fmt.Errorf("RAG_V2_MANIFEST_PRODUCTION: %w", err)
		}
	}
	return nil
}
