package model

type App struct {
	Filename      string
	Vendor        string
	To            string
	Tag           []Tag
	Addresses     []Address
	AddressGroups []AddressGroup
	Policies      []Policy
}
