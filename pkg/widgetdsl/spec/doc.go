// Package spec defines the typed intent model for the hard-cutover Widget DSL
// v2 work. These structs sit between Goja builder APIs and the existing React
// Widget IR renderer: authors build specs, specs validate intent, and lowering
// code converts validated specs into serializable Widget IR nodes.
package spec
