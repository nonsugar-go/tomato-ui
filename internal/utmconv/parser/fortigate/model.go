package fortigate

type FortiGateConfig struct {
	FirewallAddress       []map[string]AddressDetail               `yaml:"firewall_address,omitempty"`
	FirewallAddrgrp       []map[string]FirewallAddrgrpDetail       `yaml:"firewall_addrgrp,omitempty"`
	FirewallServiceCustom []map[string]FirewallServiceCustomDetail `yaml:"firewall_service_custom,omitempty"`
	FirewallServiceGroup  []map[string]FirewallServiceGroupDetail  `yaml:"firewall_service_group,omitempty"`
	FirewallPolicy        []map[string]FirewallPolicyDetail        `yaml:"firewall_policy,omitempty"`
}

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

type FirewallAddrgrpDetail struct {
	Name   string   `yaml:"-"`
	UUID   string   `yaml:"uuid"`
	Member []string `yaml:"member"`
}

type FirewallServiceCustomDetail struct {
	Name           string `yaml:"-"`
	Category       string `yaml:"category,omitempty"`
	Protocol       string `yaml:"protocol,omitempty"`
	ProtocolNumber string `yaml:"protocol-number,omitempty"`
	Icmptype       string `yaml:"icmptype,omitempty"`
	Icmpcode       string `yaml:"icmpcode,omitempty"`
	TcpPortrange   string `yaml:"tcp-portrange,omitempty"`
	UdpPortrange   string `yaml:"udp-portrange,omitempty"`
	Proxy          string `yaml:"proxy,omitempty"`
}

type FirewallServiceGroupDetail struct {
	Name   string   `yaml:"-"`
	Member []string `yaml:"member"`
}

type FirewallPolicyDetail struct {
	No                  string   `yaml:"-"`
	UUID                string   `yaml:"uuid"`
	Name                string   `yaml:"name,omitempty"`
	Srcintf             string   `yaml:"srcintf,omitempty"`
	Dstintf             string   `yaml:"dstintf,omitempty"`
	Action              string   `yaml:"action,omitempty"`
	Srcaddr             []string `yaml:"srcaddr,omitempty"`
	Dstaddr             []string `yaml:"dstaddr,omitempty"`
	InternetService     string   `yaml:"internet-service,omitempty"`
	InternetServiceName []string `yaml:"internet-service-name,omitempty"`
	Schedule            string   `yaml:"schedule,omitempty"`
	Service             []string `yaml:"service,omitempty"`
	UtmStatus           string   `yaml:"utm-status,omitempty"`
	InspectionMode      string   `yaml:"inspection-mode,omitempty"`
	SslSshProfile       string   `yaml:"ssl-ssh-profile,omitempty"`
	AvProfile           string   `yaml:"av-profile,omitempty"`
	WebfilterProfile    string   `yaml:"webfilter-profile,omitempty"`
	DnsfilterProfile    string   `yaml:"dnsfilter-profile,omitempty"`
	ApplicationList     string   `yaml:"application-list,omitempty"`
	Logtraffic          string   `yaml:"logtraffic,omitempty"`
	Nat                 string   `yaml:"nat,omitempty"`
	MatchVip            string   `yaml:"match-vip,omitempty"`
	Comments            string   `yaml:"comments,omitempty"`
}
