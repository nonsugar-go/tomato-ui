package model

type App struct {
	Filename      string
	Vendor        string
	To            string
	Addresses     []Address
	AddressGroups []AddressGroup
}
