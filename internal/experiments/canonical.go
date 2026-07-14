// Package experiments provides immutable experiment-artifact identity helpers.
package experiments

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
)

// CanonicalJSON marshals JSON-compatible values with recursively sorted object
// keys. Arrays preserve their declared order. Callers must normalize
// set-valued fields before calling this function.
func CanonicalJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, errors.Wrap(err, "marshal canonical JSON input")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, errors.Wrap(err, "decode canonical JSON input")
	}
	var output bytes.Buffer
	if err := encode(&output, decoded); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Fingerprint returns a schema-namespaced SHA-256 identifier over canonical
// JSON. Schema versions are semantic input and must change for an incompatible
// meaning change rather than mutating a published artifact.
func Fingerprint(schema string, value any) (string, error) {
	if schema == "" {
		return "", errors.New("fingerprint schema is required")
	}
	canonical, err := CanonicalJSON(value)
	if err != nil {
		return "", err
	}
	digest := sha256.New()
	_, _ = digest.Write([]byte(schema))
	_, _ = digest.Write([]byte{0})
	_, _ = digest.Write(canonical)
	return "sha256:" + hex.EncodeToString(digest.Sum(nil)), nil
}

func encode(output *bytes.Buffer, value any) error {
	switch typed := value.(type) {
	case nil, bool, string, json.Number:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return errors.Wrap(err, "encode canonical scalar")
		}
		output.Write(encoded)
		return nil
	case []any:
		output.WriteByte('[')
		for index, element := range typed {
			if index > 0 {
				output.WriteByte(',')
			}
			if err := encode(output, element); err != nil {
				return err
			}
		}
		output.WriteByte(']')
		return nil
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		output.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				output.WriteByte(',')
			}
			encodedKey, err := json.Marshal(key)
			if err != nil {
				return errors.Wrap(err, "encode canonical object key")
			}
			output.Write(encodedKey)
			output.WriteByte(':')
			if err := encode(output, typed[key]); err != nil {
				return err
			}
		}
		output.WriteByte('}')
		return nil
	default:
		return errors.Errorf("unsupported canonical JSON value %T", value)
	}
}
