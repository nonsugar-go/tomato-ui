package fortigate

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/ui"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/exporter/excel"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
	"go.yaml.in/yaml/v4"
)

func parseYAML(filename string) (*FortiGateConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := FgValidYaml(file)
	if err != nil {
		return nil, err
	}

	ui.Info("デバッグ用の yaml ファイルを出力します: fg_conf_valid.yaml")
	if err := os.WriteFile("fg_conf_valid.yaml", data, 0o644); err != nil {
		log.Fatal(err)
	}

	var config FortiGateConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	for _, entry := range config.FirewallAddress {
		for key := range entry {
			detail := entry[key]
			detail.Name = key
			entry[key] = detail
		}
	}

	for _, entry := range config.FirewallAddrgrp {
		for key := range entry {
			detail := entry[key]
			detail.Name = key
			entry[key] = detail
		}
	}

	for _, entry := range config.FirewallServiceCustom {
		for key := range entry {
			detail := entry[key]
			detail.Name = key
			entry[key] = detail
		}
	}

	for _, entry := range config.FirewallServiceGroup {
		for key := range entry {
			detail := entry[key]
			detail.Name = key
			entry[key] = detail
		}
	}

	for _, entry := range config.FirewallPolicy {
		for key := range entry {
			detail := entry[key]
			detail.No = key
			entry[key] = detail
		}
	}

	return &config, nil
}

func ParseFortiGate(app *model.App) {
	ui.Warn("FortiGate の解析は作成中です")

	var config *FortiGateConfig
	var err error

	if config, err = parseYAML(app.Filename); err != nil {
		slog.Error("Failed to parse YAML", "error", err)
		os.Exit(1)
	}

	e := excel.NewExcel("fortigate_param.xlsx")
	defer e.Close()

	e.NewSheet("Addresses")
	e.Println("Name", "UUID", "Type", "SubType", "AssociatedInterface",
		"Subnet", "Fqdn", "StartIp", "EndIp", "Dirty", "Comment")
	for _, element := range config.FirewallAddress {
		for _, v := range element {
			e.Println(v.Name, v.UUID, v.Type, v.SubType, v.AssociatedInterface,
				v.Subnet, v.Fqdn, v.StartIp, v.EndIp, v.Dirty, v.Comment)
		}
	}
	e.AddTable()

	e.NewSheet("Addrgrp")
	e.Println("Name", "UUID", "member")
	for _, element := range config.FirewallAddrgrp {
		for _, v := range element {
			e.Println(v.Name, v.UUID, strings.Join(v.Member, ";"))
		}
	}
	e.AddTable()

	e.NewSheet("ServiceCustom")
	e.Println("Name", "Category", "Protocol", "ProtocolNumber", "Icmptype", "Icmpcode",
		"TcpPortrange", "UdpPortrange", "Proxy")
	for _, element := range config.FirewallServiceCustom {
		for _, v := range element {
			e.Println(v.Name, v.Category, v.Protocol, v.ProtocolNumber, v.Icmptype, v.Icmpcode,
				v.TcpPortrange, v.UdpPortrange, v.Proxy)
		}
	}
	e.AddTable()

	e.NewSheet("ServiceGroup")
	e.Println("Name", "Member")
	for _, element := range config.FirewallServiceGroup {
		for _, v := range element {
			e.Println(v.Name, strings.Join(v.Member, ";"))
		}
	}
	e.AddTable()

	e.NewSheet("Policy")
	e.Println("No", "UUID", "Name", "Srcintf", "Dstintf", "Action",
		"Srcaddr", "Dstaddr",
		"InternetService", "InternetServiceName",
		"Schedule", "Service", "UtmStatus", "InspectionMode",
		"SslSshProfile", "AvProfile", "WebfilterProfile", "DnsfilterProfile", "ApplicationList",
		"Logtraffic", "Nat", "MatchVip", "Comments")
	for _, element := range config.FirewallPolicy {
		for _, v := range element {
			e.Println(v.No, v.UUID, v.Name, v.Srcintf, v.Dstintf, v.Action,
				strings.Join(v.Srcaddr, ";"), strings.Join(v.Dstaddr, ";"),
				v.InternetService, strings.Join(v.InternetServiceName, ";"),
				v.Schedule, strings.Join(v.Service, ";"), v.UtmStatus, v.InspectionMode,
				v.SslSshProfile, v.AvProfile, v.WebfilterProfile, v.DnsfilterProfile, v.ApplicationList,
				v.Logtraffic, v.Nat, v.MatchVip, v.Comments)
		}
	}
	e.AddTable()
}
