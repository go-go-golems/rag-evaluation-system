package ragcontract

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

type Issue struct {
	Code     string   `json:"code"`
	Path     string   `json:"path"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Operator string   `json:"operator,omitempty"`
	Hint     string   `json:"hint,omitempty"`
}
type Report struct {
	Issues []Issue `json:"issues"`
}

func (r Report) OK() bool {
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			return false
		}
	}
	return true
}
func (r *Report) Add(code, path, message string) {
	r.Issues = append(r.Issues, Issue{Code: code, Path: path, Message: message, Severity: SeverityError})
}
func (r *Report) Normalize() {
	sort.Slice(r.Issues, func(i, j int) bool {
		if r.Issues[i].Path == r.Issues[j].Path {
			return r.Issues[i].Code < r.Issues[j].Code
		}
		return r.Issues[i].Path < r.Issues[j].Path
	})
}

type ValidationError struct{ Report Report }

func (e *ValidationError) Error() string {
	values := make([]string, 0, len(e.Report.Issues))
	for _, i := range e.Report.Issues {
		if i.Severity == SeverityError {
			values = append(values, fmt.Sprintf("%s at %s: %s", i.Code, i.Path, i.Message))
		}
	}
	return strings.Join(values, "; ")
}

var identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,127}$`)

func ValidIdentifier(value string) bool { return identifierPattern.MatchString(value) }

func ValidatePipeline(value PipelineIR) Report {
	r := Report{}
	if value.SchemaVersion != PipelineSchemaVersion {
		r.Add("RAG_V2_SCHEMA_VERSION", "$.schemaVersion", "expected "+PipelineSchemaVersion)
	}
	inputIDs := map[string]bool{}
	for i, v := range value.Inputs {
		p := fmt.Sprintf("$.inputs[%d]", i)
		if !ValidIdentifier(v.ID) {
			r.Add("RAG_V2_INPUT_ID", p+".id", "invalid input ID")
		}
		if inputIDs[v.ID] {
			r.Add("RAG_V2_DUPLICATE_INPUT", p+".id", "duplicate input ID")
		}
		inputIDs[v.ID] = true
		if v.Kind == "" {
			r.Add("RAG_V2_PORT_KIND", p+".kind", "port kind is required")
		}
		if v.BindingMode != "artifact" && v.BindingMode != "request" && v.BindingMode != "dataset" {
			r.Add("RAG_V2_INPUT_BINDING_MODE", p+".bindingMode", "binding mode must be artifact, request, or dataset")
		}
		if v.BindingMode == "artifact" && (v.ArtifactRole == "" || v.ManifestSchema == "") {
			r.Add("RAG_V2_INPUT_MANIFEST", p, "artifact role and manifest schema are required")
		}
	}
	nodes := map[string]Node{}
	for i, n := range value.Nodes {
		p := fmt.Sprintf("$.nodes[%d]", i)
		if !ValidIdentifier(n.ID) {
			r.Add("RAG_V2_NODE_ID", p+".id", "invalid node ID")
		}
		if _, ok := nodes[n.ID]; ok {
			r.Add("RAG_V2_DUPLICATE_NODE", p+".id", "duplicate node ID")
		}
		nodes[n.ID] = n
		if !ValidIdentifier(n.Operator.Kind) || !regexp.MustCompile(`^v[1-9][0-9]*$`).MatchString(n.Operator.Version) {
			r.Add("RAG_V2_OPERATOR_REF", p+".operator", "operator must use a namespaced kind and vN version")
		}
		if _, err := CanonicalRaw(n.Config, "{}"); err != nil {
			r.Add("RAG_V2_CONFIG_JSON", p+".config", err.Error())
		}
		ports := map[string]bool{}
		for j, b := range n.Inputs {
			bp := fmt.Sprintf("%s.inputs[%d]", p, j)
			if ports[b.Port] {
				r.Add("RAG_V2_DUPLICATE_PORT", bp+".port", "duplicate input port")
			}
			ports[b.Port] = true
		}
	}
	for i, n := range value.Nodes {
		for j, b := range n.Inputs {
			if _, ok := nodes[b.From.NodeID]; !ok && !inputIDs[b.From.NodeID] {
				r.Add("RAG_V2_UNKNOWN_REFERENCE", fmt.Sprintf("$.nodes[%d].inputs[%d].from.nodeId", i, j), "unknown node or input")
			}
		}
	}
	for i, o := range value.Outputs {
		if _, ok := nodes[o.From.NodeID]; !ok && !inputIDs[o.From.NodeID] {
			r.Add("RAG_V2_UNKNOWN_REFERENCE", fmt.Sprintf("$.outputs[%d].from.nodeId", i), "unknown node or input")
		}
	}
	r.Normalize()
	return r
}

func RequireValid(report Report) error {
	if report.OK() {
		return nil
	}
	return &ValidationError{Report: report}
}
