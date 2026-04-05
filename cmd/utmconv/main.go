package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/converter/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/paloalto"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func writeLines(filename string, lines []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		fmt.Fprintln(f, line)
	}
	return nil
}

func main() {
	handler := log.New(os.Stderr)
	handler.SetLevel(log.DebugLevel)
	handler.SetReportTimestamp(true)
	slog.SetDefault(slog.New(handler))
	var app model.App

	flag.StringVar(&app.Filename, "in", "", "comfig file")
	flag.StringVar(&app.Utm, "utm", "panorama", "utm type")
	flag.StringVar(&app.To, "to", "cp", "output format")
	flag.Parse()

	var confirm bool = false
	if app.Filename != "" && app.Utm != "" {
		confirm = true
	}

	for {
		if confirm {
			break
		}
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewFilePicker().
					Title("ファイル名").
					Description("Panorama の xml ファイルを選択してください").
					CurrentDirectory(".").
					DirAllowed(true).
					// AllowedTypes([]string{".xml"}).
					Value(&app.Filename),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("解析する UTM の種類").
					Options(
						huh.NewOption("Panorama", "panorama"),
						huh.NewOption("PaloAlto", "pa"),
						huh.NewOption("FortiGate", "fg"),
						huh.NewOption("Checkpoint", "cp"),
					).
					Value(&app.Utm),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		slog.Info("設定ファイル", "config_file", app.Filename)
		slog.Info("対象", "utm", app.Utm)

		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("入力は正しいですか？").
					Affirmative("はい").
					Negative("いいえ").
					Value(&confirm),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}
	}

	switch app.Utm {
	case "panorama":
		paloalto.ParsePanorama(&app)
		slog.Info("Panorama の解析が終了しました", "output", "Panorama.xlsx")
		switch app.To {
		case "json":
			slog.Warn("JSON output is not implemented yet")
		case "cp":
			lines, err := checkpoint.ConvertAddresses(app.Addresses)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeLines("checkpoint_address.conf", lines)
			slog.Info("Check Point のアドレス変換が終了しました",
				"output", "checkpoint_address.conf")

			lines, err = checkpoint.ConvertAddressGroups(app.AddressGroups)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeLines("checkpoint_address_group.conf", lines)
			slog.Info("Check Point のアドレスグループ変換が終了しました",
				"output", "checkpoint_address_group.conf")

		default:
			slog.Error("unsupported output", "to", app.To)
		}
	default:
		slog.Error("UTM の指定は未実装です", "utm", app.Utm)
	}
}
