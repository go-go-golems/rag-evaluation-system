package spec

// JSONValue is the serialized value domain accepted by Widget IR and action
// payloads. Builders may carry richer Go values internally, but public Widget IR
// output must reduce to this shape before crossing the browser boundary.
type JSONValue any

// JSONObject is a convenience alias for JSON object payloads.
type JSONObject map[string]JSONValue

// PageSpec is the typed v2 authoring model for a Widget page.
type PageSpec struct {
	SchemaVersion string
	ID            string
	Title         string
	Meta          JSONObject
	Root          NodeSpec
	Diagnostics   []ValidationIssue
}

// NodeSpec is the common typed representation used before lowering to the
// current Widget IR node maps.
type NodeSpec struct {
	Kind        NodeKind
	Type        string
	Text        string
	Tag         string
	Props       JSONObject
	Children    []NodeSpec
	Source      *SourceSpan
	Diagnostics []ValidationIssue
}

type NodeKind string

const (
	NodeKindText      NodeKind = "text"
	NodeKindElement   NodeKind = "element"
	NodeKindComponent NodeKind = "component"
)

// SourceSpan records authoring provenance for diagnostics and future debug
// overlays. It is intentionally optional at the renderer boundary.
type SourceSpan struct {
	Module string
	Helper string
	Path   string
}

// SectionSpec captures intent-level document sectioning before it becomes a
// SectionBlock component node.
type SectionSpec struct {
	Title    string
	Level    int
	AnchorID string
	Caption  string
	Children []NodeSpec
}

// SchemaSpec defines ordered record fields. The ordered slice is deliberate:
// JavaScript object insertion order is not the contract v2 should depend on.
type SchemaSpec struct {
	Name   string
	Fields []FieldSpec
}

// FieldSpec describes one record field as a set of explicit facets rather than
// the v1 single role string doing storage, semantic, layout, and editor work.
type FieldSpec struct {
	Name       string
	Label      string
	Kind       FieldKind
	Semantic   FieldSemantic
	Editor     EditorSpec
	Summary    SummarySpec
	Layout     FieldLayout
	Validation FieldValidation
}

type FieldKind string

const (
	FieldKindString  FieldKind = "string"
	FieldKindNumber  FieldKind = "number"
	FieldKindBoolean FieldKind = "boolean"
	FieldKindDate    FieldKind = "date"
	FieldKindMedia   FieldKind = "media"
	FieldKindURL     FieldKind = "url"
)

type FieldSemantic string

const (
	FieldSemanticKey     FieldSemantic = "key"
	FieldSemanticPrimary FieldSemantic = "primary"
	FieldSemanticShort   FieldSemantic = "short"
	FieldSemanticProse   FieldSemantic = "prose"
	FieldSemanticCount   FieldSemantic = "count"
	FieldSemanticSize    FieldSemantic = "size"
	FieldSemanticMeasure FieldSemantic = "measure"
	FieldSemanticStatus  FieldSemantic = "status"
	FieldSemanticTags    FieldSemantic = "tags"
)

type EditorSpec struct {
	Control     EditorControl
	Placeholder string
	Rows        int
	ReadOnly    bool
}

type EditorControl string

const (
	EditorControlText     EditorControl = "text"
	EditorControlTextarea EditorControl = "textarea"
	EditorControlSelect   EditorControl = "select"
	EditorControlCheckbox EditorControl = "checkbox"
)

type SummarySpec struct {
	CellKind string
	Elide    bool
}

type FieldLayout struct {
	Width string
}

type FieldValidation struct {
	Required  bool
	MinLength int
	MaxLength int
}

// CollectionSpec is the v2 intent model for tables, selectable tables,
// master-detail editors, and later richer multi-view collections.
type CollectionSpec struct {
	Name        string
	Rows        []JSONObject
	Schema      SchemaSpec
	Mode        CollectionMode
	Selection   *SelectionSpec
	Arrangement ArrangementSpec
	Actions     CollectionActions
	Empty       string
}

type CollectionMode string

const (
	CollectionModeShow   CollectionMode = "show"
	CollectionModeEdit   CollectionMode = "edit"
	CollectionModePick   CollectionMode = "pick"
	CollectionModeManage CollectionMode = "manage"
)

type SelectionSpec struct {
	Kind  SelectionKind
	Param string
	Value string
}

type SelectionKind string

const (
	SelectionKindURLParam SelectionKind = "urlParam"
)

type ArrangementSpec struct {
	Kind ArrangementKind
}

type ArrangementKind string

const (
	ArrangementKindTable        ArrangementKind = "table"
	ArrangementKindMasterDetail ArrangementKind = "master-detail"
)

type CollectionActions struct {
	Open    *ActionSpec
	Create  *CreateActionSpec
	Submit  *SubmitSpec
	Reorder *ActionSpec
	Remove  *ActionSpec
}

type CreateActionSpec struct {
	Label string
}

type SubmitSpec struct {
	FormAction string
	Method     string
}

// ActionSpec describes browser-visible behavior as data. JavaScript callbacks
// may configure builders or register server handlers, but the browser receives
// this serializable shape.
type ActionSpec struct {
	Kind    ActionKind
	Name    string
	To      string
	Event   string
	Payload PayloadTemplate
	Confirm *TemplateSpec
	Result  *ServerResultPolicy
}

type ActionKind string

const (
	ActionKindNavigate ActionKind = "navigate"
	ActionKindServer   ActionKind = "server"
	ActionKindDownload ActionKind = "download"
	ActionKindEvent    ActionKind = "event"
	ActionKindCopy     ActionKind = "copy"
)

type PayloadTemplate struct {
	Fields []PayloadFieldSpec
}

type PayloadFieldSpec struct {
	Name  string
	Value TemplateValue
}

type TemplateSpec struct {
	Parts []TemplateValue
}

type TemplateValue struct {
	Kind  TemplateValueKind
	Path  string
	Value JSONValue
	Text  string
}

type TemplateValueKind string

const (
	TemplateValueLiteral TemplateValueKind = "literal"
	TemplateValuePath    TemplateValueKind = "path"
	TemplateValueText    TemplateValueKind = "text"
)

type ServerResultPolicy struct {
	Refresh bool
	Toast   bool
	Patch   bool
}

// ValidationIssue is the common diagnostic unit returned by specs, builders,
// and future page validation endpoints.
type ValidationIssue struct {
	Severity ValidationSeverity
	Code     string
	Path     string
	Message  string
	Hint     string
}

type ValidationSeverity string

const (
	ValidationSeverityInfo    ValidationSeverity = "info"
	ValidationSeverityWarning ValidationSeverity = "warning"
	ValidationSeverityError   ValidationSeverity = "error"
)
