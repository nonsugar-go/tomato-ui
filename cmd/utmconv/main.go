package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

type tickMsg time.Time

type modelTUI struct {
	spinner   spinner.Model
	loading   bool
	quitting  bool
	fadeStep  int
	fadingOut bool
	mode      string
}

func initialModel() modelTUI {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	return modelTUI{
		spinner: s,
		loading: true,
		mode:    "splash",
	}
}

func (m modelTUI) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			time.Sleep(1 * time.Second)
			return tickMsg{}
		},
	)
}

func (m modelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var spinnerColors = []string{"196", "202", "208", "214"}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		idx := int(time.Now().UnixNano()/1e8) % len(spinnerColors)
		m.spinner.Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(spinnerColors[idx]))

		return m, cmd

	case tickMsg:
		m.fadingOut = true
		return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			return t
		})

	case time.Time:
		if m.fadingOut {
			m.fadeStep++

			if m.fadeStep > 10 {
				m.fadingOut = false
				m.mode = "header"
				return m, tea.Quit
			}

			return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
				return t
			})
		}
	}

	return m, nil
}

func (m modelTUI) View() string {
	if m.quitting {
		return ""
	}

	switch m.mode {
	case "splash":
		if m.fadingOut {
			return m.fadeView()
		}
		return m.splashView()

	case "header":
		return m.headerView()
	}

	return ""
}

var logoText string = `
	_____                          _             _   _ ___ 
	|_   _|__  _ __ ___   __ _ _ __| |_ ___      | | | |_ _|
	| |/ _ \| '_ ' _ \ / _' | '__| __/ _ \_____| | | || | 
	| | (_) | | | | | | (_| | |  | || (_) |_____| |_| || | 
	|_|\___/|_| |_| |_|\__,_|_|   \__\___/      \___/|___|`

var utmText string = `
	_   _ _ __ ___   ___ ___  _ __ __   __
	| | | | '_ ' _ \ / __/ _ \| '_ \\ \ / /
	| |_| | | | | | | (_| (_) | | | |\ V / 
	\___/|_| |_| |_|\___\___/|_| |_| \_/  
	`

func (m modelTUI) splashView() string {
	colors := []string{"196", "202", "208", "214", "220"}

	logoLines := strings.Split(strings.Trim("\n"+logoText, "\n"), "\n")
	var coloredLogo []string
	for i, line := range logoLines {
		color := colors[i%len(colors)]
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Align(lipgloss.Center)
		coloredLogo = append(coloredLogo, style.Render(line))
	}
	logoBlock := strings.Join(coloredLogo, "\n")

	logoWidth := lipgloss.Width(logoBlock)

	utmStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Align(lipgloss.Center).
		Width(logoWidth)

	utmBlock := utmStyle.Render(utmText)

	tomatoLine := lipgloss.NewStyle().
		Width(logoWidth).
		Align(lipgloss.Center).
		Render(strings.Repeat("🍅 ", logoWidth/4))

	loadingMsg := lipgloss.NewStyle().
		Width(logoWidth).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%s  Initializing Tomato Systems...", m.spinner.View()))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoBlock,
		utmBlock,
		tomatoLine,
		loadingMsg,
	)
	return "\n\n" + content + "\n\n"
}

func (m modelTUI) fadeView() string {
	fade := 1.0 - float64(m.fadeStep)/10.0
	if fade < 0 {
		fade = 0
	}

	color := lipgloss.Color(fmt.Sprintf("%d", 232+int(fade*20)))

	style := lipgloss.NewStyle().
		Foreground(color).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		style.Render(logoText),
		style.Render(utmText),
		style.Render(strings.Repeat("🍅 ", 60/4)),
		style.Render(fmt.Sprintf("%s Initializing...", m.spinner.View())),
	)

	return "\n\n" + content + "\n\n"
}

func (m modelTUI) headerView() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Render("Tomato-UI")

	sub := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("utmconv")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		sub,
	) + "\n\n🍅 Ready!\n\n"
}

func writeMgmtLines(filename string, lines []string, app model.App) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, `mgmt_cli login -u secadmin -p Lab@12345 >sid.txt`)
	for _, line := range lines {
		fmt.Fprint(f, "mgmt_cli "+line)
		if app.IgnoreWarnings {
			fmt.Fprint(f, ` ignore-warnings true`)
		}
		fmt.Fprintln(f, ` -s sid.txt`)
	}
	fmt.Fprintln(f, `mgmt_cli discard -s sid.txt
# mgmt_cli publish -s sid.txt
mgmt_cli logout -s sid.txt
rm sid.txt`)
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
	flag.BoolVar(&app.IgnoreWarnings, "ignore-warnings", false, "ignore warnings during conversion")
	flag.Parse()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

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
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("解析する UTM のベンダーを選択してください").
					Description("⏎ を押した後に、ベンダーを選択してください\n\n\n\n\n\n\n").
					Options(
						huh.NewOption("Panorama", "panorama"),
						// huh.NewOption("PaloAlto", "pa"),
						// huh.NewOption("FortiGate", "fg"),
						// huh.NewOption("Checkpoint", "cp"),
					).
					Value(&app.Vendor),
			),
			huh.NewGroup(
				huh.NewFilePicker().
					Title("設定ファイルを選択してください").
					Description("⏎ を押した後に、ファイルを選択してください").
					CurrentDirectory(".").
					DirAllowed(true).
					// AllowedTypes([]string{".xml"}).
					Value(&app.Filename),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("変換形式を選択してください").
					Options(
						huh.NewOption("変換しない", ""),
						huh.NewOption("Check Point CLI", "cp"),
						// huh.NewOption("JSON", "json"),
					).
					Value(&app.To),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("変換時の警告を無視しますか？").
					Affirmative("はい").
					Negative("いいえ").
					Value(&app.IgnoreWarnings),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		slog.Info("対象ベンダ", "vendor", app.Vendor)
		slog.Info("設定ファイル", "config_file", app.Filename)
		slog.Info("変換形式", "to", app.To)
		slog.Info("変換時の警告を無視", "ignore_warnings", app.IgnoreWarnings)

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
		slog.Info("Panorama の解析が終了しました", "output", "panorama.xlsx")
		switch app.To {
		case "":
			slog.Info("変換しないが選択されました。処理を終了します")
		case "json":
			slog.Warn("JSON output is not implemented yet")
		case "cp":
			ctx := checkpoint.NewContext()
			lines, err := checkpoint.ConvertAddresses(app.Addresses, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_address.conf", lines, app)
			slog.Info("Check Point のアドレス変換が終了しました",
				"output", "checkpoint_address.conf")

			lines, err = checkpoint.ConvertAddressGroups(app.AddressGroups, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_address_group.conf", lines, app)
			slog.Info("Check Point のアドレスグループ変換が終了しました",
				"output", "checkpoint_address_group.conf")

			lines, err = checkpoint.ConvertServices(app.Services)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_service.conf", lines, app)
			slog.Info("Check Point のサービス変換が終了しました",
				"output", "checkpoint_service.conf")

			lines, err = checkpoint.ConvertServiceGroups(app.ServiceGroups)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_service_group.conf", lines, app)
			slog.Info("Check Point のサービスグループ変換が終了しました",
				"output", "checkpoint_service_group.conf")

			lines, err = checkpoint.ConvertPolicies(app.Policies, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_policy.conf", lines, app) // NOTE: untested
			slog.Info("Check Point のポリシー変換が終了しました",
				"output", "checkpoint_policy.conf")

		default:
			slog.Error("unsupported output", "to", app.To)
		}
	default:
		slog.Error("Vendor の指定は未実装です", "vendor", app.Vendor)
	}

	if interactive {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("処理が完了しました。⏎ キーを押して終了してください").
					Value(&confirm),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}
	}

	slog.Info("終了します")
}
