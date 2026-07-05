package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/ui"
	conv_checkpoint "github.com/nonsugar-go/tomato-ui/internal/utmconv/converter/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/fortigate"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/paloalto"
	"github.com/pelletier/go-toml/v2"
)

const appConfigFilename = "utmconv.toml"

func loadOrInitConfig(path string) (*model.AppConfig, error) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		defaultCfg := model.NewDefaultAppConfig()

		data, err := toml.Marshal(defaultCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write config file: %w", err)
		}

		ui.Info("設定ファイルが見つからないため、デフォルト値で生成します: %s", path)
		ui.Info("設定ファイルを確認・編集後、utmconv を再実行してください")
		os.Exit(0)
	}

	ui.Info("設定ファイルを読み込みます: %s", path)
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg model.AppConfig

	if err := toml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &cfg, nil
}

func writeMgmtLines(filename string, lines []string, app model.App) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "mgmt_cli login -u %s -p %s >sid.txt\n",
		app.AppConfig.CheckPoint.Cli.MgmtCliUser.Value,
		app.AppConfig.CheckPoint.Cli.MgmtCliPassword.Value)
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			fmt.Fprintln(f, line)
		} else {
			fmt.Fprint(f, "mgmt_cli "+line)
			if app.AppConfig.CheckPoint.Cli.IgnoreWarnings.Value {
				fmt.Fprint(f, ` ignore-warnings true`)
			}
			fmt.Fprintln(f, ` -s sid.txt`)
		}
	}
	fmt.Fprintln(f, `#
mgmt_cli discard -s sid.txt
# mgmt_cli publish -s sid.txt
mgmt_cli logout -s sid.txt
rm sid.txt`)
	return nil
}

func main() {
	fmt.Println(`🥫🍅 Tomato-UI
🍅🥫 utmconv v1.0.0
----------------------------------------`)

	var app model.App

	cfg, err := loadOrInitConfig(appConfigFilename)
	if err != nil {
		ui.Error("設定ファイルの読み込みに失敗: %s: %s", appConfigFilename, err.Error())
		os.Exit(1)
	}

	app.AppConfig = *cfg

	flag.StringVar(&app.Filename, "in", "", "config file")
	flag.StringVar(&app.Vendor, "vendor", "", "vendor type")
	flag.StringVar(&app.To, "to", "", "output format")
	flag.Parse()

	var confirm bool = false
	var interactive bool = true
	if app.Filename != "" && app.Vendor != "" {
		confirm = true
		interactive = false
	}

	for {
		if confirm {
			break
		}
		fmt.Println("対応している UTM (ベンダー)")
		fmt.Println("  - Check Point: cp")
		fmt.Println("  - Panorama: panorama")
		fmt.Println("  - FortiGate (作成中): fg")
		for {
			input := ui.Prompt("解析する UTM を選択してください (cp/panorama/fg)", "cp")
			switch input {
			case "cp", "panorama", "fg":
				app.Vendor = input
			case "":
				app.Vendor = "cp"
			default:
				app.Vendor = ""
			}
			if app.Vendor != "" {
				break
			}
		}

		fmt.Println("設定ファイルを入力してください")
		globPattern := "*"
		switch app.Vendor {
		case "cp":
			fmt.Println("\"show_package-YYYY-MM-DD_HH-MM-SS.tar.gz\" を指定してください")
			globPattern = "*.tar.gz"
		case "panorama":
			fmt.Println("Panorama から取り出した xml ファイルを指定してください")
			globPattern = "*.xml"
		case "fg":
			fmt.Println("FortiGate から取り出した yaml ファイルを指定してください")
			globPattern = "*.yaml"
		}
		app.Filename = ui.SelectFile("設定ファイルを入力してください", app.Filename, globPattern)

		fmt.Println("解析後に他の UTM 用のコンフィグを出力する場合は、変換対象を指定してください")
		fmt.Println("対応している UTM (ベンダー)")
		fmt.Println("  - Check Point: cp")
		fmt.Println("  - Panorama: panorama")
		fmt.Println("  - FortiGate: fg")
		for {
			input := ui.Prompt("変換対象を指定してください", "変換しない")
			switch input {
			case "cp", "panorama", "fg":
				app.To = input
			case "変換しない":
				app.To = ""
			default:
				app.Vendor = "invalid"
			}
			if app.Vendor != "invalid" {
				break
			}
		}

		fmt.Println("----------------------------------------")
		ui.Info("対象ベンダ: %s", app.Vendor)
		ui.Info("設定ファイル: %s", app.Filename)
		ui.Info("変換形式: %s", app.To)
		fmt.Println("----------------------------------------")

		confirm = ui.Confirm("入力は正しいですか？", true)
	}

	switch app.Vendor {
	case "panorama":
		paloalto.ParsePanorama(&app)
		ui.Success("Panorama の解析が終了しました: %s", "panorama_param.xlsx")
		switch app.To {
		case "":
			ui.Info("変換しないが選択されました。処理を終了します")
		case "json":
			ui.Warn("JSON output is not implemented yet")
		case "cp":
			ctx := conv_checkpoint.NewContext(&app)
			lines, err := conv_checkpoint.ConvertAddresses(app.Addresses, ctx)
			if err != nil {
				ui.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_address.conf", lines, app)
			ui.Success("Check Point のアドレス変換が終了しました: %s", "checkpoint_address.conf")

			lines, err = conv_checkpoint.ConvertAddressGroups(app.AddressGroups, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_address_group.conf", lines, app)
			ui.Success("Check Point のアドレスグループ変換が終了しました: %s", "checkpoint_address_group.conf")

			lines, err = conv_checkpoint.ConvertServices(app.Services, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_service.conf", lines, app)
			ui.Success("Check Point のサービス変換が終了しました: %s", "checkpoint_service.conf")

			lines, err = conv_checkpoint.ConvertServiceGroups(app.ServiceGroups, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_service_group.conf", lines, app)
			ui.Success("Check Point のサービスグループ変換が終了しました: %s", "checkpoint_service_group.conf")

			lines, err = conv_checkpoint.ConvertPolicies(app.Policies, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_policy.conf", lines, app)
			ui.Success("Check Point のポリシー変換が終了しました: %s", "checkpoint_policy.conf")

			lines, err = conv_checkpoint.ConvertNATPolicies(app.NATRules, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_nat.conf", lines, app)
			ui.Success("Check Point の NAT 変換が終了しました: %s", "checkpoint_nat.conf")

			// 脅威ポリシーの変換は、app.Policies をもとに行う。
			lines, err = conv_checkpoint.ConvertThreatPolicies(app.Policies, ctx)
			if err != nil {
				ui.Error("convert error: %s", err.Error())
			}
			writeMgmtLines("checkpoint_threat.conf", lines, app)
			ui.Success("Check Point の脅威防御の変換が終了しました: %s", "checkpoint_threat.conf")

		default:
			ui.Error("unsupported output: %s", app.To)
		}

	case "cp":
		checkpoint.ParseCheckPoint(&app)
		ui.Success("Check Point の解析が終了しました: %s", "checkpoint_param.xlsx")
		switch app.To {
		case "":
			ui.Info("変換しないが選択されました。処理を終了します")
		default:
			ui.Error("unsupported output: %s", app.To)
		}

	case "fg":
		fortigate.ParseFortiGate(&app)
		ui.Success("FortiGate の解析が終了しました: %s", "fortigate_param.xlsx")
		switch app.To {
		case "":
			ui.Info("変換しないが選択されました。処理を終了します")
		default:
			ui.Error("unsupported output: %s", app.To)
		}

	default:
		ui.Error("Vendor の指定は未実装です: %s", app.Vendor)
	}

	if interactive {
		_ = ui.Confirm("⏎ キーを押して終了してください", true)
		// form := huh.NewForm(
		// 	huh.NewGroup(
		// 		huh.NewConfirm().
		// 			Title("処理が完了しました。⏎ キーを押して終了してください").
		// 			Value(&confirm),
		// 	),
		// )
		// if err := form.Run(); err != nil {
		// 	log.Fatal(err)
		// }
	}

	ui.Info("終了します")
}
