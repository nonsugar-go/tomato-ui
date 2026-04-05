package model

type App struct {
	Filename      string
	Utm           string
	To            string
	Addresses     []Address
	AddressGroups []AddressGroup
}
