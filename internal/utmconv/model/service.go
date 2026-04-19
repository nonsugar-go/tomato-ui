package model

type ServiceType string

const (
	ServiceTypeUnknown ServiceType = "unknown"
	ServiceTypeTCP     ServiceType = "tcp"
	ServiceTypeUDP     ServiceType = "udp"
	ServiceTypeSCTP    ServiceType = "sctp"
	ServiceTypeICMP    ServiceType = "icmp"
	ServiceTypeICMPv6  ServiceType = "icmpv6"
	ServiceTypeOther   ServiceType = "other"
)

type Service struct {
	Name        string
	Type        ServiceType
	Ports       string // "80", "20-30", "any"
	SourcePorts string // "80", "20-30", "any"
	Description string
	Tags        []string
}

type ServiceGroup struct {
	Name        string
	Members     []string // Service.Name or ServiceGroup.Name
	Description string
	Tags        []string
}
