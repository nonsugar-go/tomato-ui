package fortigate

// TODO:
type FortiGateConfig struct {
	FirewallAddress []map[string]AddressDetail `yaml:"firewall_address,omitempty"`
}

// TODO:
type AddressDetail struct {
	Name                string `yaml:"-"`
	UUID                string `yaml:"uuid"`
	Comment             string `yaml:"comment,omitempty"`
	Type                string `yaml:"type,omitempty"`
	AssociatedInterface string `yaml:"associated-interface,omitempty"`
	Subnet              string `yaml:"subnet,omitempty"`
	Fqdn                string `yaml:"fqdn,omitempty"`
	StartIp             string `yaml:"start-ip,omitempty"`
	EndIp               string `yaml:"end-ip,omitempty"`
	SubType             string `yaml:"sub-type,omitempty"`
	Dirty               string `yaml:"dirty,omitempty"`
}
