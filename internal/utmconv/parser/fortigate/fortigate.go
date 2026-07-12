package fortigate

import (
	"log"
	"log/slog"
	"os"

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

	slog.Info("デバッグ用の yaml ファイルを出力します。", slog.String("file", "fg_conf_valid.yaml"))
	if err := os.WriteFile("fg_conf_valid.yaml", data, 0o644); err != nil {
		log.Fatal(err)
	}

	var config FortiGateConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	for _, entry := range config.FirewallAddress {
		for key, detail := range entry {
			detail.Name = key
		}
	}

	return &config, nil
}

func ParseFortiGate(app *model.App) {
	slog.Error("FortiGate の解析は作成中です。", "vendor", app.To)

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
		for k, v := range element {
			e.Println(k, v.UUID, v.Type, v.SubType, v.AssociatedInterface,
				v.Subnet, v.Fqdn, v.StartIp, v.EndIp, v.Dirty, v.Comment)
		}
	}
	e.AddTable()
}
