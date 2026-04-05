package model

type AddressType string

const (
	AddressTypeUnknown   AddressType = "unknown"
	AddressTypeIPNetmask AddressType = "ip-netmask"
	AddressTypeHost      AddressType = "host"
	AddressTypeNetwork   AddressType = "network"
	AddressTypeRange     AddressType = "range"
	AddressTypeFQDN      AddressType = "fqdn"
)

type Address struct {
	Name        string
	Type        AddressType
	Value       string
	Description string
	Tags        []string
}

type AddressGroup struct {
	Name        string
	Members     []string // Address.Name or AddressGroup.Name
	Description string
	Tags        []string
}
