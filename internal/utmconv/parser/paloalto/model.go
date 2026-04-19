package paloalto

type PaloAltoConfig struct {
	TagObject     []ScopedTagObject
	Addresses     []ScopedAddress
	AddressGroups []ScopedAddressGroup
	Services      []ScopedService
	ServiceGroups []ScopedServiceGroup
	SecurityRules []ScopedSecurity
	NATRules      []ScopedNAT
}

type ScopedTagObject struct {
	Scope     string // shared / dg
	TagObject TagObject
}

type ScopedAddress struct {
	Scope   string // shared / dg
	Address Address
}

type ScopedAddressGroup struct {
	Scope string // shared / dg
	Group AddressGroup
}

type ScopedService struct {
	Scope   string // shared / dg
	Service Service
}

type ScopedServiceGroup struct {
	Scope string // shared / dg
	Group ServiceGroup
}

type ScopedSecurity struct {
	Scope        string // shared / dg
	Rulebase     string // pre / post
	SecurityRule SecurityRule
}

type ScopedNAT struct {
	Scope    string // shared / dg
	Rulebase string // pre / post
	NATRule  NATRule
}

type Config struct {
	Shared   Shared   `xml:"shared"`
	Devices  Devices  `xml:"devices"`
	Policies Policies `xml:"policies"`
}

type Shared struct {
	Tags              []TagObject        `xml:"tag>entry"` // NOTE: untested
	Addresses         []Address          `xml:"address>entry"`
	AddressGroups     []AddressGroup     `xml:"address-group>entry"`
	Services          []Service          `xml:"service>entry"`
	ServiceGroups     []ServiceGroup     `xml:"service-group>entry"`
	Applications      []Application      `xml:"application>entry"`
	ApplicationGroups []ApplicationGroup `xml:"application-group>entry"`
	PreRulebase       Rulebase           `xml:"pre-rulebase"`
	PostRulebase      Rulebase           `xml:"post-rulebase"`
}

type Devices struct {
	DeviceGroups []DeviceGroup `xml:"entry>device-group>entry"`
}

type DeviceGroup struct {
	Name              string             `xml:"name,attr"`
	Tags              []TagObject        `xml:"tag>entry"` // NOTE: untested
	Addresses         []Address          `xml:"address>entry"`
	AddressGroups     []AddressGroup     `xml:"address-group>entry"`
	Services          []Service          `xml:"service>entry"`
	ServiceGroups     []ServiceGroup     `xml:"service-group>entry"`
	Applications      []Application      `xml:"application>entry"`
	ApplicationGroups []ApplicationGroup `xml:"application-group>entry"`
	PreRulebase       Rulebase           `xml:"pre-rulebase"`
	PostRulebase      Rulebase           `xml:"post-rulebase"`
}

type Rulebase struct {
	Security Security `xml:"security"`
	Nat      NAT      `xml:"nat"`
}

type Policies struct {
	Security Security `xml:"security"`
}

type TagObject struct {
	Name     string `xml:"name,attr"`
	Color    string `xml:"color"`
	Comments string `xml:"comments"` // NOTE: untested
}

type Address struct {
	Name        string   `xml:"name,attr"`
	IPNetmask   string   `xml:"ip-netmask"`
	IPRange     string   `xml:"ip-range"`    // NOTE: untested
	IPWildcard  string   `xml:"ip-wildcard"` // NOTE: untested
	FQDN        string   `xml:"fqdn"`
	Description string   `xml:"description"`
	Tags        []string `xml:"tag>member"`
}

type AddressGroup struct {
	Name        string   `xml:"name,attr"`
	Static      []string `xml:"static>member"`
	Dynamic     *Dynamic `xml:"dynamic"` // NOTE: untested
	Description string   `xml:"description"`
	Tags        []string `xml:"tag>member"` // NOTE: untested
}

type Dynamic struct {
	Filter string `xml:"filter"` // NOTE: untested
}

type Service struct {
	Name        string   `xml:"name,attr"`
	Protocol    Protocol `xml:"protocol"`
	Description string   `xml:"description"`
	Tags        []string `xml:"tag>member"` // NOTE: untested
}

type Protocol struct {
	TCP *TCP `xml:"tcp"`
	UDP *UDP `xml:"udp"`
}

type TCP struct {
	Port       string `xml:"port"`
	SourcePort string `xml:"source-port"` // NOTE: untested
}

type UDP struct {
	Port       string `xml:"port"`
	SourcePort string `xml:"source-port"`
}

type Security struct {
	Rules []SecurityRule `xml:"rules>entry"`
}

type ServiceGroup struct {
	Name        string   `xml:"name,attr"`
	Members     []string `xml:"members>member"`
	Description string   `xml:"description"` // NOTE: untested
	Tags        []string `xml:"tag>member"`  // NOTE: untested
}

type Application struct {
	Name        string   `xml:"name,attr"`
	Category    string   `xml:"category"`    // NOTE: untested
	Subcategory string   `xml:"subcategory"` // NOTE: untested
	Technology  string   `xml:"technology"`  // NOTE: untested
	Risk        string   `xml:"risk"`        // NOTE: untested
	Default     *Default `xml:"default"`     // NOTE: untested
	Description string   `xml:"description"` // NOTE: untested
	Tags        []string `xml:"tag>member"`  // NOTE: untested
}

type Default struct {
	Port string `xml:"port"` // NOTE: untested
}

type ApplicationGroup struct {
	Name        string   `xml:"name,attr"`
	Members     []string `xml:"members>member"`
	Description string   `xml:"description"`
	Tags        []string `xml:"tag>member"`
}

type SecurityRule struct {
	Name     string   `xml:"name,attr"`
	Tags     []string `xml:"tag>member"`
	GroupTag string   `xml:"group-tag"`

	FromZones    []string `xml:"from>member"` // NOTE: untested
	ToZones      []string `xml:"to>member"`   // NOTE: untested
	Sources      []string `xml:"source>member"`
	Destinations []string `xml:"destination>member"`
	Applications []string `xml:"application>member"`
	Services     []string `xml:"service>member"`
	Action       string   `xml:"action"`

	Disabled          string `xml:"disabled"`
	NegateSource      string `xml:"negate-source"`
	NegateDestination string `xml:"negate-destination"` // NOTE: untested
	LogStart          string `xml:"log-start"`
	LogEnd            string `xml:"log-end"`
	LogSetting        string `xml:"log-setting"`

	Schedule string `xml:"schedule"` // NOTE: untested

	SourceUsers    []string `xml:"source-user>member"`     // NOTE: untested
	Categories     []string `xml:"category>member"`        // NOTE: untested
	SourceHIP      []string `xml:"source-hip>member"`      // NOTE: untested
	DestinationHIP []string `xml:"destination-hip>member"` // NOTE: untested

	ProfileSetting *ProfileSetting `xml:"profile-setting"` // NOTE: untested
	Target         *Target         `xml:"target"`          // NOTE: untested

	Description string `xml:"description"`
}

type ProfileSetting struct {
	Group []string `xml:"group>member"`

	Profiles *Profiles `xml:"profiles"` // NOTE: untested
}

type Profiles struct {
	AV  []string `xml:"virus>member"`             // Antivirus
	VP  []string `xml:"vulnerability>member"`     // Vulnerability Protection
	AS  []string `xml:"spyware>member"`           // Anti-Spyware
	URL []string `xml:"url-filtering>member"`     // URL Filtering
	FB  []string `xml:"file-blocking>member"`     // File Blocking
	DF  []string `xml:"data-filtering>member"`    // Data Filtering
	WFA []string `xml:"wildfire-analysis>member"` // WildFire Analysis
}

type Target struct {
	Negate  string        `xml:"negate"`        // NOTE: untested
	Devices []TargetEntry `xml:"devices>entry"` // NOTE: untested
}

type TargetEntry struct {
	Name string `xml:"name,attr"` // NOTE: untested
}

type NAT struct {
	Rules []NATRule `xml:"rules>entry"`
}

type NATRule struct {
	Name string `xml:"name,attr"`

	UUID     string   `xml:"uuid"`
	Disabled string   `xml:"disabled"` // yes/no
	Tags     []string `xml:"tag>member"`

	FromZones    []string `xml:"from>member"`
	ToZones      []string `xml:"to>member"`
	Sources      []string `xml:"source>member"`
	Destinations []string `xml:"destination>member"`
	Service      string   `xml:"service"`

	ToInterfaces []string `xml:"to-interface>member"`
	SourceUsers  []string `xml:"source-user>member"`
	Category     []string `xml:"category>member"`

	// Source NAT
	SourceTranslation *SourceTranslation `xml:"source-translation"`

	// Destination NAT
	DestinationTranslation *DestinationTranslation `xml:"destination-translation"`

	Description string `xml:"description"`
}

type SourceTranslation struct {
	DynamicIPAndPort *DynamicIPAndPort `xml:"dynamic-ip-and-port"`
	StaticIP         *StaticIP         `xml:"static-ip"`
}

type DynamicIPAndPort struct {
	TranslatedAddress []string          `xml:"translated-address>member"`
	InterfaceAddress  *InterfaceAddress `xml:"interface-address"`

	AddressType string `xml:"address-type"` // interface-address / translated-address
}

type InterfaceAddress struct {
	Interface  string `xml:"interface"`
	Ip         string `xml:"ip"`
	FloatingIp string `xml:"floating-ip"`
}

type StaticIP struct {
	TranslatedAddress string `xml:"translated-address"`
	BiDirectional     string `xml:"bi-directional"` // yes/no
}

type DestinationTranslation struct {
	TranslatedAddress string `xml:"translated-address"`
	TranslatedPort    string `xml:"translated-port"`

	Protocol string `xml:"protocol"`
}

// colorMap maps Palo Alto's color names to more human-readable color names.
// Note: The actual color names used by Palo Alto may differ,
//
//	and this mapping is based on common color naming conventions.
//	Adjust as necessary based on the specific color codes used in your Palo Alto configuration.
var colorMap = map[string]string{
	"color1":  "Red",
	"color2":  "Green",
	"color3":  "Blue",
	"color4":  "Yellow",
	"color5":  "Copper",
	"color6":  "Orange",
	"color7":  "Purple",
	"color8":  "Gray",
	"color9":  "Light Green",
	"color10": "Cyan",
	"color11": "Light Gray",
	"color12": "Blue Gray",
	"color13": "Lime",
	"color14": "Black",
	"color15": "Gold",

	"color16": "Brown",
	"color17": "Olive",
	"color18": "Maroon",
	"color19": "Navy",
	"color20": "Teal",
	"color21": "Aqua",
	"color22": "Fuchsia",
	"color23": "Silver",
	"color24": "Dark Gray",
	"color25": "Light Blue",
	"color26": "Light Yellow",
	"color27": "Light Cyan",
	"color28": "Light Pink",
	"color29": "Light Orange",
	"color30": "Light Purple",

	"color31": "Dark Red",
	"color32": "Dark Green",
	"color33": "Dark Blue",
	"color34": "Dark Yellow",
	"color35": "Dark Orange",
	"color36": "Dark Purple",
	"color37": "Dark Cyan",
	"color38": "Dark Brown",
	"color39": "Dark Olive",
	"color40": "Dark Teal",
	"color41": "Dark Gray 2",
	"color42": "Dark Black",
}
