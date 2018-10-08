// Package pkgsrc in order to build macOS package from source
package pkgsrc

// Pkg definition
type Pkg struct {
	Name     string
	Vers     string
	Ext      string
	namesrc  string
	URL      string
	HashType string
	Hash     string
}
