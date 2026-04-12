package model

type Tag struct {
	Key         string // optional: category (FortiGate), empty for Palo Alto
	Value       string // required
	Color       string // optional: display only (non-portable)
	Description string // optional
}
