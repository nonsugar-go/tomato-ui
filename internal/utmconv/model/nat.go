package model

type NATType string

const (
	NATTypeStatic  NATType = "static"  // Fixed 1:1 address translation
	NATTypeDynamic NATType = "dynamic" // Dynamic 1:1 translation using an address pool
	NATTypeHide    NATType = "hide"    // Many-to-one translation with port overloading (PAT)
	NATTypePAT     NATType = "pat"     // Port Translation
	NATTypeNoNAT   NATType = "no-nat"  // NAT exemption / identity NAT

)

type NATRule struct {
	ID string

	Name        string
	Enabled     bool
	Description string
	Tags        []string

	FromZones []string
	ToZones   []string

	OriginalSource      []string
	OriginalDestination []string
	OriginalService     []string

	TranslatedSource      []string
	TranslatedDestination []string
	TranslatedService     []string

	Type NATType

	// Applies only to static NAT
	BiDirectional bool

	// Palo Alto: shared, dg:<name>, pre, post
	// Check Point: install-on target
	Scope string
}
