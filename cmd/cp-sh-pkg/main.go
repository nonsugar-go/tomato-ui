package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/exporter/excel"
)

// Object
type CPObject struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
	Type string `json:"type"`

	Comments string     `json:"comments,omitempty"`
	Tags     []CPObject `json:"tags,omitempty"`
	Color    string     `json:"color,omitempty"`

	// host
	IPv4Address string `json:"ipv4-address,omitempty"`
	IPv6Address string `json:"ipv6-address,omitempty"`

	// network
	Subnet4     string `json:"subnet4,omitempty"`
	MaskLength4 int    `json:"mask-length4,omitempty"`
	Subnet6     string `json:"subnet6,omitempty"`
	MaskLength6 int    `json:"mask-length6"`

	// dns-domain
	IsSubDomain bool `json:"is-sub-domain,omitempty"`

	// group / service-group
	Members []string `json:"members,omitempty"`

	// service-tcp / service-udp
	Port       string `json:"port,omitempty"`
	SourcePort string `json:"source-port,omitempty"`
}

// Access Rule
type AccessRule struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
	Type string `json:"type"`

	// "access-section"
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`

	// "access-rule"
	Comments string     `json:"comments,omitempty"`
	Tags     []CPObject `json:"tags,omitempty"`

	RuleNumber        int      `json:"rule-number,omitempty"`
	Enabled           bool     `json:"enabled,omitempty"`
	SourceNegate      bool     `json:"source-negate,omitempty"`
	Source            []string `json:"source,omitempty"`
	DestinationNegate bool     `json:"destination-negate,omitempty"`
	Destination       []string `json:"destination,omitempty"`
	ServiceNegate     bool     `json:"service-negate,omitempty"`
	Service           []string `json:"service,omitempty"`
	ServiceResource   string   `json:"service-resource,omitempty"`
	ContentDirection  string   `json:"content-direction,omitempty"`
	ContentNegate     bool     `json:"content-negate,omitempty"`
	Content           []string `json:"content,omitempty"`

	Vpn    []string `json:"vpn,omitempty"`
	Action string   `json:"action,omitempty"`

	Track struct {
		PerSession            bool   `json:"per-session,omitempty"`
		PerConnection         bool   `json:"per-connection,omitempty"`
		Alert                 string `json:"alert,omitempty"`
		EnableFirewallSession bool   `json:"enable-firewall-session,omitempty"`
		Accounting            bool   `json:"accounting,omitempty"`
		Type                  string `json:"type,omitempty"`
	} `json:"track"`

	InstallOn    []string `json:"install-on,omitempty"`
	CustomFields struct {
		Field1 string `json:"field-1,omitempty"`
		Field2 string `json:"field-2,omitempty"`
		Field3 string `json:"field-3,omitempty"`
	} `json:"custom-fields"`
}

func joinNames(objs []CPObject) string {
	var sb strings.Builder
	for i, o := range objs {
		sb.WriteString(o.Name)
		if i != 0 {
			sb.WriteRune(';')
		}
	}
	return sb.String()
}

func parseAccessRules(e *excel.Excel, inDir string) {
	pattern := filepath.Join(inDir, "* Network-Management server.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("fail to find files",
			slog.String("pattern", pattern), slog.String("error", err.Error()))
		os.Exit(1)
	}

	for _, file := range files {
		ruleName := strings.TrimSuffix(filepath.Base(file), " Network-Management server.json")
		var accessRules []AccessRule
		jsonData, err := os.ReadFile(file)
		if err != nil {
			slog.Error("Faild to read a file", slog.String("filename", file), slog.String("error", err.Error()))
			continue
		}
		err = json.Unmarshal(jsonData, &accessRules)
		if err != nil {
			slog.Error("JSON error", slog.String("error", err.Error()))
			continue
		}

		e.NewSheet("Access RUles - " + ruleName)
		e.Println("No", "Name", "Enabled", "Comments")
		for _, r := range accessRules {
			e.Println(strconv.Itoa(r.RuleNumber), r.Name, fmt.Sprintf("%t", r.Enabled),
				r.Comments)
		}
		e.AddTable()
	}
}

func parseCPObjects(e *excel.Excel, inDir string) {
	pattern := filepath.Join(inDir, "*_objects.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("fail to find files",
			slog.String("pattern", pattern), slog.String("error", err.Error()))
		os.Exit(1)
	}

	var tmpCPObjs []CPObject
	for _, file := range files {
		slog.Info("processing", slog.String("file", file))

		data, err := os.ReadFile(file)
		if err != nil {
			slog.Error("fail to read file",
				slog.String("file", file), slog.String("error", err.Error()))
			continue
		}

		var objs []CPObject
		if err := json.Unmarshal(data, &objs); err != nil {
			slog.Error("fail to unmarshal JSON",
				slog.String("file", file), slog.String("error", err.Error()))
			continue
		}
		tmpCPObjs = append(tmpCPObjs, objs...)
	}

	exists := make(map[string]struct{}, len(tmpCPObjs))
	cpObjs := make([]CPObject, 0, len(tmpCPObjs))
	for _, obj := range tmpCPObjs {
		if _, ok := exists[obj.Uid]; ok {
			continue
		}
		exists[obj.Uid] = struct{}{}
		cpObjs = append(cpObjs, obj)
	}

	slog.Info("count of objects", slog.Int("count", len(cpObjs)))

	byType := make(map[string][]*CPObject)
	for i := range cpObjs {
		obj := &cpObjs[i]
		byType[obj.Type] = append(byType[obj.Type], obj)
	}
	for _, objs := range byType {
		sort.Slice(objs, func(i, j int) bool {
			a := strings.ToLower(objs[i].Name)
			b := strings.ToLower(objs[j].Name)
			if a == b {
				return objs[i].Name < objs[j].Name
			}
			return a < b
		})
	}

	keys := make([]string, 0, len(byType))
	for k := range byType {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a := strings.ToLower(keys[i])
		b := strings.ToLower(keys[j])
		if a == b {
			return keys[i] < keys[j]
		}
		return a < b
	})

	for _, typeName := range keys {
		count := len(byType[typeName])
		slog.Info("objects of type", slog.String("type", typeName), slog.Int("count", count))
	}

	mapUidObj := make(map[string]*CPObject, len(cpObjs))
	for i := range cpObjs {
		obj := &cpObjs[i]
		mapUidObj[obj.Uid] = obj
	}

	writeCPObjectsExcel(e, byType, mapUidObj)
}

func writeCPObjectsExcel(e *excel.Excel, byType map[string][]*CPObject, mapUidObj map[string]*CPObject) error {
	e.NewSheet("host")
	e.Println("Name", "IPv4 Addr", "IPv6 Addr", "Tags", "Comments")
	for _, o := range byType["host"] {
		e.Println(o.Name, o.IPv4Address, o.IPv6Address, joinNames(o.Tags), o.Comments)
	}
	e.AddTable()

	e.NewSheet("network")
	e.Println("Name", "Subnet4", "len4", "Subnet6", "len6", "Tags", "Comments")
	for _, o := range byType["network"] {
		e.Println(o.Name, o.Subnet4, strconv.Itoa(o.MaskLength4), o.Subnet6, strconv.Itoa(o.MaskLength6),
			joinNames(o.Tags), o.Comments)
	}
	e.AddTable()

	e.NewSheet("dns-domain")
	e.Println("Name", "Is-sub-domain", "Tags", "Comments")
	for _, o := range byType["dns-domain"] {
		e.Println(o.Name, fmt.Sprintf("%t", o.IsSubDomain), joinNames(o.Tags), o.Comments)
	}
	e.AddTable()

	e.NewSheet("group")
	e.Println("Name", "members", "Tags", "Comments")
	for _, o := range byType["group"] {
		members := make([]string, 0, len(o.Members))
		for _, uid := range o.Members {
			members = append(members, mapUidObj[uid].Name)
		}
		sort.Slice(members, func(i, j int) bool {
			a := strings.ToLower(members[i])
			b := strings.ToLower(members[j])
			if a == b {
				return members[i] < members[j]
			}
			return a < b
		})

		if len(members) == 0 {
			e.Println(o.Name, "", joinNames(o.Tags), o.Comments)
		} else {
			for i, m := range members {
				if i == 0 {
					e.Println(o.Name, m, joinNames(o.Tags), o.Comments)
				} else {
					e.Println("", m)
				}
			}
		}
	}
	e.AddTable()

	e.NewSheet("service-tcp")
	e.Println("Name", "port", "source-port", "Tags", "Comments")
	for _, o := range byType["service-tcp"] {
		e.Println(o.Name, o.Port, o.SourcePort, joinNames(o.Tags), o.Comments)
	}
	e.AddTable()

	e.NewSheet("service-udp")
	e.Println("Name", "port", "source-port", "Tags", "Comments")
	for _, o := range byType["service-udp"] {
		e.Println(o.Name, o.Port, o.SourcePort, joinNames(o.Tags), o.Comments)
	}
	e.AddTable()

	e.NewSheet("service-group")
	e.Println("Name", "members", "Tags", "Comments")
	for _, o := range byType["service-group"] {
		members := make([]string, 0, len(o.Members))
		for _, uid := range o.Members {
			members = append(members, mapUidObj[uid].Name)
		}
		sort.Slice(members, func(i, j int) bool {
			a := strings.ToLower(members[i])
			b := strings.ToLower(members[j])
			if a == b {
				return members[i] < members[j]
			}
			return a < b
		})

		if len(members) == 0 {
			e.Println(o.Name, "", joinNames(o.Tags), o.Comments)
		} else {
			for i, m := range members {
				if i == 0 {
					e.Println(o.Name, m, joinNames(o.Tags), o.Comments)
				} else {
					e.Println("", m)
				}
			}
		}
	}
	e.AddTable()

	return nil
}

func main() {
	handler := log.New(os.Stderr)
	handler.SetLevel(log.DebugLevel)
	handler.SetReportTimestamp(true)
	slog.SetDefault(slog.New(handler))

	var inDir string
	var outFile string
	flag.StringVar(&inDir, "dir", "./", "show package json files directory")
	flag.StringVar(&outFile, "out", "check_point_param.xlsx", "output Excel file name")

	flag.Parse()

	slog.Info("json files directory", slog.String("dir", inDir))
	slog.Info("output file", slog.String("file", outFile))

	e := excel.NewExcel(outFile)
	defer e.Close()
	parseCPObjects(e, inDir)
	parseAccessRules(e, inDir)
	slog.Info("Excel file saved", slog.String("file", outFile))
}
