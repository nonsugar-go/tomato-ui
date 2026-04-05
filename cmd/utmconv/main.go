package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/converter/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/paloalto"

	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// メッセージの定義
type tickMsg time.Time

type modelTUI struct {
	spinner  spinner.Model
	loading  bool
	quitting bool
}

func initialModel() modelTUI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("204")) // トマト色
	return modelTUI{spinner: s, loading: true}
}

func (m modelTUI) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		// 2秒後に終了するコマンド（実際はここで初期化処理を行う）
		func() tea.Msg {
			time.Sleep(2 * time.Second)
			return tickMsg{}
		},
	)
}

func (m modelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tickMsg:
		m.loading = false
		return m, tea.Quit // 起動完了。本来はここでメイン画面のModelへ切り替える
	}
	return m, nil
}

func (m modelTUI) View() string {
	if m.quitting {
		return ""
	}

	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6347")).Padding(0, 1)
	tomato := "🍅"

	logoText := `
  _____                          _             _   _ ___ 
 |_   _|__  _ __ ___   __ _ _ __| |_ ___      | | | |_ _|
   | |/ _ \| '_ ' _ \ / _' | '__| __/ _ \_____| | | || | 
   | | (_) | | | | | | (_| | |  | || (_) |_____| |_| || | 
   |_|\___/|_| |_| |_|\__,_|_|   \__\___/      \___/|___|`

	// 上下のトマト列を作る
	topRow := lipgloss.NewStyle().MarginLeft(4).Render("🍅 🍅 🍅 🍅 🍅 🍅 🍅 🍅 🍅")
	bottomRow := lipgloss.NewStyle().MarginLeft(4).Render("🍅 🍅 🍅 🍅 🍅 🍅 🍅 🍅 🍅")

	// ロゴの左右にトマトを配置
	middleRow := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tomato,
		logoStyle.Render(logoText),
		tomato,
	)

	// スピナーとテキストの結合
	loadingMsg := fmt.Sprintf("\n  %s  Initializing Tomato Systems...", m.spinner.View())

	// 全体を垂直に結合
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		topRow,
		middleRow,
		bottomRow,
		loadingMsg,
	)

	return "\n\n" + content + "\n\n"
}

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
	flag.StringVar(&app.Vendor, "vendor", "panorama", "vendor type")
	flag.StringVar(&app.To, "to", "cp", "output format")
	flag.Parse()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// スプラッシュ終了後の処理
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Render("✓ Ready! Welcome to tomato-ui."))

	var confirm bool = false
	if app.Filename != "" && app.Vendor != "" {
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
					Title("解析する UTM のベンダー").
					Options(
						huh.NewOption("Panorama", "panorama"),
						huh.NewOption("PaloAlto", "pa"),
						huh.NewOption("FortiGate", "fg"),
						huh.NewOption("Checkpoint", "cp"),
					).
					Value(&app.Vendor),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("変換形式").
					Options(
						huh.NewOption("変換しない", ""),
						huh.NewOption("Check Point CLI", "cp"),
						huh.NewOption("JSON", "json"),
					).
					Value(&app.To),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		slog.Info("設定ファイル", "config_file", app.Filename)
		slog.Info("対象", "utm", app.Vendor)
		slog.Info("変換形式", "to", app.To)

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

	switch app.Vendor {
	case "panorama":
		paloalto.ParsePanorama(&app)
		slog.Info("Panorama の解析が終了しました", "output", "Panorama.xlsx")
		switch app.To {
		case "":
			slog.Info("変換しないが選択されました。処理を終了します")
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
		slog.Error("Vendor の指定は未実装です", "vendor", app.Vendor)
	}
}
