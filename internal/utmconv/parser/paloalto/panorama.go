package paloalto

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/exporter/excel"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

type PaloAltoConfig struct {
	Addresses     []ScopedAddress
	AddressGroups []ScopedAddressGroup
	Services      []ScopedService
	SecurityRules []ScopedSecurity
	NATRules      []ScopedNAT
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

type Address struct {
	Name        string `xml:"name,attr"`
	IPNetmask   string `xml:"ip-netmask"`
	FQDN        string `xml:"fqdn"`
	Description string `xml:"description"`
}

type AddressGroup struct {
	Name    string   `xml:"name,attr"`
	Static  []string `xml:"static>member"`
	Dynamic *Dynamic `xml:"dynamic"`
}

type Dynamic struct {
	Filter string `xml:"filter"`
}

type Service struct {
	Name     string   `xml:"name,attr"`
	Protocol Protocol `xml:"protocol"`
}

type Protocol struct {
	TCP *TCP `xml:"tcp"`
	UDP *UDP `xml:"udp"`
}

type TCP struct {
	Port string `xml:"port"`
}

type UDP struct {
	Port string `xml:"port"`
}

type Security struct {
	Rules []SecurityRule `xml:"rules>entry"`
}

type ServiceGroup struct {
	Name    string   `xml:"name,attr"`
	Members []string `xml:"members>member"`
}

type Application struct {
	Name        string `xml:"name,attr"`
	Category    string `xml:"category"`
	Subcategory string `xml:"subcategory"`
	Technology  string `xml:"technology"`
	Risk        string `xml:"risk"`
	Description string `xml:"description"`
}

type ApplicationGroup struct {
	Name    string   `xml:"name,attr"`
	Members []string `xml:"members>member"`
}

type SecurityRule struct {
	Name         string   `xml:"name,attr"`
	FromZones    []string `xml:"from>member"`
	ToZones      []string `xml:"to>member"`
	Sources      []string `xml:"source>member"`
	Destinations []string `xml:"destination>member"`
	Applications []string `xml:"application>member"`
	Services     []string `xml:"service>member"`
	Action       string   `xml:"action"`
	Description  string   `xml:"description"`
	Tags         []string `xml:"tag>member"`
}

type NAT struct {
	Rules []NATRule `xml:"rules>entry"`
}

type NATRule struct {
	Name         string   `xml:"name,attr"`
	FromZones    []string `xml:"from>member"`
	ToZones      []string `xml:"to>member"`
	Sources      []string `xml:"source>member"`
	Destinations []string `xml:"destination>member"`
	Service      string   `xml:"service"`

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
	TranslatedAddress []string `xml:"translated-address>member"`
}

type StaticIP struct {
	TranslatedAddress string `xml:"translated-address"`
}

type DestinationTranslation struct {
	TranslatedAddress string `xml:"translated-address"`
	TranslatedPort    string `xml:"translated-port"`
}

func BuildPaloAltoConfig(c *Config) *PaloAltoConfig {
	var result PaloAltoConfig

	for _, addr := range c.Shared.Addresses {
		result.Addresses = append(result.Addresses, ScopedAddress{
			Scope:   "shared",
			Address: addr,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, addr := range dg.Addresses {
			result.Addresses = append(result.Addresses, ScopedAddress{
				Scope:   dg.Name,
				Address: addr,
			})
		}
	}

	for _, addrGrp := range c.Shared.AddressGroups {
		result.AddressGroups = append(result.AddressGroups, ScopedAddressGroup{
			Scope: "shared",
			Group: addrGrp,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, addrGrp := range dg.AddressGroups {
			result.AddressGroups = append(result.AddressGroups, ScopedAddressGroup{
				Scope: dg.Name,
				Group: addrGrp,
			})
		}
	}

	for _, svc := range c.Shared.Services {
		result.Services = append(result.Services, ScopedService{
			Scope:   "shared",
			Service: svc,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, svc := range dg.Services {
			result.Services = append(result.Services, ScopedService{
				Scope:   dg.Name,
				Service: svc,
			})
		}
	}

	return &result
}

func printAddresses(scope string, addrs []Address, e *excel.Excel) {
	for _, addr := range addrs {
		e.Println(
			scope,
			addr.Name,
			addr.IPNetmask,
			addr.FQDN,
			addr.Description,
		)
	}
}

func printAddressGroups(scope string, groups []AddressGroup, e *excel.Excel) {
	for _, g := range groups {
		staticMembers := fmt.Sprintf("%v", g.Static)
		dynamicFilter := ""
		if g.Dynamic != nil {
			dynamicFilter = g.Dynamic.Filter
		}
		e.Println(
			scope,
			g.Name,
			staticMembers,
			dynamicFilter,
		)
	}
}

func printServices(scope string, services []Service, e *excel.Excel) {
	for _, svc := range services {
		proto := ""
		port := ""
		if svc.Protocol.TCP != nil {
			proto = "tcp"
			port = svc.Protocol.TCP.Port
		}
		if svc.Protocol.UDP != nil {
			proto = "udp"
			port = svc.Protocol.UDP.Port
		}
		e.Println(
			scope,
			svc.Name,
			proto,
			port,
		)
	}
}

func printServiceGroups(scope string, groups []ServiceGroup, e *excel.Excel) {
	for _, g := range groups {
		e.Println(
			scope,
			g.Name,
			fmt.Sprintf("%v", g.Members),
		)
	}
}

func printApplications(scope string, apps []Application, e *excel.Excel) {
	for _, app := range apps {
		e.Println(
			scope,
			app.Name,
			app.Category,
			app.Subcategory,
			app.Technology,
			app.Risk,
			app.Description,
		)
	}
}

func printApplicationGroups(scope string, groups []ApplicationGroup, e *excel.Excel) {
	for _, g := range groups {
		e.Println(
			scope,
			g.Name,
			fmt.Sprintf("%v", g.Members),
		)
	}
}

func printSecurityRules(scope string, rules []SecurityRule, e *excel.Excel) {
	for _, rule := range rules {
		e.Println(
			scope,
			rule.Name,
			fmt.Sprintf("%v", rule.FromZones),
			fmt.Sprintf("%v", rule.ToZones),
			fmt.Sprintf("%v", rule.Sources),
			fmt.Sprintf("%v", rule.Destinations),
			fmt.Sprintf("%v", rule.Applications),
			fmt.Sprintf("%v", rule.Services),
			rule.Action,
			fmt.Sprintf("%v", rule.Tags),
			rule.Description,
		)
	}
}

func printNatRules(scope string, rules []NATRule, e *excel.Excel) {
	for _, rule := range rules {

		srcTrans := ""
		if rule.SourceTranslation != nil {
			if rule.SourceTranslation.DynamicIPAndPort != nil {
				srcTrans = fmt.Sprintf("DIPP:%v",
					rule.SourceTranslation.DynamicIPAndPort.TranslatedAddress)
			}
			if rule.SourceTranslation.StaticIP != nil {
				srcTrans = "Static:" + rule.SourceTranslation.StaticIP.TranslatedAddress
			}
		}

		dstTrans := ""
		if rule.DestinationTranslation != nil {
			dstTrans = fmt.Sprintf("%s:%s",
				rule.DestinationTranslation.TranslatedAddress,
				rule.DestinationTranslation.TranslatedPort)
		}

		e.Println(
			scope,
			rule.Name,
			fmt.Sprintf("%v", rule.FromZones),
			fmt.Sprintf("%v", rule.ToZones),
			fmt.Sprintf("%v", rule.Sources),
			fmt.Sprintf("%v", rule.Destinations),
			rule.Service,
			srcTrans,
			dstTrans,
			rule.Description,
		)
	}
}

func parseXML(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ParsePanorama(app *model.App) {
	var config *Config
	var err error

	if config, err = parseXML(app.Filename); err != nil {
		log.Fatal(err)
	}

	panoramaConfig := BuildPaloAltoConfig(config)
	fillModel(app, panoramaConfig)

	e := excel.NewExcel("panorama.xlsx")
	defer e.Close()

	e.NewSheet("Addresses")
	e.Println("Scope", "Name", "IP/Netmask", "FQDN", "Description")
	printAddresses("shared", config.Shared.Addresses, e)
	for _, dg := range config.Devices.DeviceGroups {
		printAddresses(dg.Name, dg.Addresses, e)
	}
	e.AddTable()

	e.NewSheet("Address Groups")
	e.Println("Scope", "Name", "Static Members", "Dynamic Filter")
	printAddressGroups("shared", config.Shared.AddressGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printAddressGroups(dg.Name, dg.AddressGroups, e)
	}
	e.AddTable()

	e.NewSheet("Services")
	e.Println("Scope", "Name", "Protocol", "Port")
	printServices("shared", config.Shared.Services, e)
	for _, dg := range config.Devices.DeviceGroups {
		printServices(dg.Name, dg.Services, e)
	}
	e.AddTable()

	e.NewSheet("Service Groups")
	e.Println("Scope", "Name", "Members")
	printServiceGroups("shared", config.Shared.ServiceGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printServiceGroups(dg.Name, dg.ServiceGroups, e)
	}
	e.AddTable()

	e.NewSheet("Applications")
	e.Println("Scope", "Name", "Category", "Subcategory", "Technology", "Risk", "Description")
	printApplications("shared", config.Shared.Applications, e)
	for _, dg := range config.Devices.DeviceGroups {
		printApplications(dg.Name, dg.Applications, e)
	}
	e.AddTable()

	e.NewSheet("Application Groups")
	e.Println("Scope", "Name", "Members")
	printApplicationGroups("shared", config.Shared.ApplicationGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printApplicationGroups(dg.Name, dg.ApplicationGroups, e)
	}
	e.AddTable()

	e.NewSheet("Security Rules")
	e.Println("名前", "From", "To", "Source", "Destination", "Application",
		"Service", "Action", "Tag", "Description")
	printSecurityRules("shared-pre", config.Shared.PreRulebase.Security.Rules, e)
	printSecurityRules("shared-post", config.Shared.PostRulebase.Security.Rules, e)
	for _, dg := range config.Devices.DeviceGroups {
		printSecurityRules(dg.Name+"-pre", dg.PreRulebase.Security.Rules, e)
		printSecurityRules(dg.Name+"-post", dg.PostRulebase.Security.Rules, e)
	}
	e.AddTable()

	e.NewSheet("NAT Rules")
	e.Println("Scope", "Name", "From", "To", "Source", "Destination",
		"Service", "Src Trans", "Dst Trans", "Description")
	printNatRules("shared-pre", config.Shared.PreRulebase.Nat.Rules, e)
	printNatRules("shared-post", config.Shared.PostRulebase.Nat.Rules, e)
	for _, dg := range config.Devices.DeviceGroups {
		printNatRules(dg.Name+"-pre", dg.PreRulebase.Nat.Rules, e)
		printNatRules(dg.Name+"-post", dg.PostRulebase.Nat.Rules, e)
	}
	e.AddTable()
}
