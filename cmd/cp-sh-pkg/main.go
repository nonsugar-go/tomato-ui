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

type CPObject struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
	Type string `json:"type"`

	Comments string   `json:"comments,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Color    string   `json:"color,omitempty"`

	// host
	IPv4Address string `json:"ipv4-address,omitempty"`
	IPv6Address string `json:"ipv6-address,omitempty"`

	// network
	Subnet4     string `json:"subnet4,omitempty"`
	MaskLength4 int32  `json:"mask-length4,omitempty"`

	// dns-domain
	IsSubDomain bool `json:"is-sub-domain,omitempty"`

	// group / service-group
	Members []string `json:"members,omitempty"`

	// service-tcp / service-udp
	Port       string `json:"port,omitempty"`
	SourcePort string `json:"source-port,omitempty"`
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

	writeExcel(outFile, byType, mapUidObj)
}

func writeExcel(fileName string, byType map[string][]*CPObject, mapUidObj map[string]*CPObject) error {
	e := excel.NewExcel(fileName)
	defer e.Close()

	e.NewSheet("host")
	e.Println("Name", "IPv4", "Tag", "Comments")
	for _, o := range byType["host"] {
		e.Println(o.Name, o.IPv4Address, strings.Join(o.Tags, ";"), o.Comments)
	}
	e.AddTable()

	e.NewSheet("network")
	e.Println("Name", "Subnet4", "Mask-length4", "Tag", "Comments")
	for _, o := range byType["network"] {
		e.Println(o.Name, o.Subnet4, strconv.Itoa(int(o.MaskLength4)), strings.Join(o.Tags, ";"), o.Comments)
	}
	e.AddTable()

	e.NewSheet("dns-domain")
	e.Println("Name", "Is-sub-domain", "Tag", "Comments")
	for _, o := range byType["dns-domain"] {
		e.Println(o.Name, fmt.Sprintf("%t", o.IsSubDomain), strings.Join(o.Tags, ";"), o.Comments)
	}
	e.AddTable()

	e.NewSheet("group")
	e.Println("Name", "members", "Tag", "Comments")
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
			e.Println(o.Name, "", strings.Join(o.Tags, ";"), o.Comments)
		} else {
			for i, m := range members {
				if i == 0 {
					e.Println(o.Name, m, strings.Join(o.Tags, ";"), o.Comments)
				} else {
					e.Println("", m)
				}
			}
		}
	}
	e.AddTable()

	e.NewSheet("service-tcp")
	e.Println("Name", "port", "source-port", "Tag", "Comments")
	for _, o := range byType["service-tcp"] {
		e.Println(o.Name, o.Port, o.SourcePort, strings.Join(o.Tags, ";"), o.Comments)
	}
	e.AddTable()

	e.NewSheet("service-udp")
	e.Println("Name", "port", "source-port", "Tag", "Comments")
	for _, o := range byType["service-udp"] {
		e.Println(o.Name, o.Port, o.SourcePort, strings.Join(o.Tags, ";"), o.Comments)
	}
	e.AddTable()

	e.NewSheet("service-group")
	e.Println("Name", "members", "Tag", "Comments")
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
			e.Println(o.Name, "", strings.Join(o.Tags, ";"), o.Comments)
		} else {
			for i, m := range members {
				if i == 0 {
					e.Println(o.Name, m, strings.Join(o.Tags, ";"), o.Comments)
				} else {
					e.Println("", m)
				}
			}
		}
	}
	e.AddTable()

	slog.Info("Excel file saved", slog.String("file", fileName))
	return nil
}
