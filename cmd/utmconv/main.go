package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	conv_checkpoint "github.com/nonsugar-go/tomato-ui/internal/utmconv/converter/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/parser/checkpoint"
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

const appConfigFilename = "utmconv.json"

func loadOrInitConfig(path string) (*model.AppConfig, error) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		defaultCfg := model.NewDefaultAppConfig()

		data, err := json.MarshalIndent(defaultCfg, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write config file: %w", err)
		}

		slog.Info("設定ファイルが見つからないため、デフォルト値で生成します", "filename", path)
		slog.Info("設定ファイルを確認・編集後、utmconv を再実行してください")
		os.Exit(0)
	}

	slog.Info("設定ファイルを読み込みます", "filename", path)
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg model.AppConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
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
	handler := log.New(os.Stderr)
	handler.SetLevel(log.DebugLevel)
	handler.SetReportTimestamp(true)
	slog.SetDefault(slog.New(handler))
	var app model.App

	cfg, err := loadOrInitConfig(appConfigFilename)
	if err != nil {
		slog.Error("設定ファイルの読み込みに失敗", "error", err.Error(), "filename", appConfigFilename)
		os.Exit(1)
	}
	app.AppConfig = *cfg

	flag.StringVar(&app.Filename, "in", "", "comfig file")
	flag.StringVar(&app.Vendor, "vendor", "panorama", "vendor type")
	flag.StringVar(&app.To, "to", "cp", "output format")
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
						huh.NewOption("Checkpoint (作成中)", "cp"),
						// huh.NewOption("PaloAlto", "pa"),
						// huh.NewOption("FortiGate", "fg"),
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
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		slog.Info("対象ベンダ", "vendor", app.Vendor)
		slog.Info("設定ファイル", "config_file", app.Filename)
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
		slog.Info("Panorama の解析が終了しました", "output", "panorama.xlsx")
		switch app.To {
		case "":
			slog.Info("変換しないが選択されました。処理を終了します")
		case "json":
			slog.Warn("JSON output is not implemented yet")
		case "cp":
			ctx := conv_checkpoint.NewContext(&app)
			lines, err := conv_checkpoint.ConvertAddresses(app.Addresses, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_address.conf", lines, app)
			slog.Info("Check Point のアドレス変換が終了しました",
				"output", "checkpoint_address.conf")

			lines, err = conv_checkpoint.ConvertAddressGroups(app.AddressGroups, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_address_group.conf", lines, app)
			slog.Info("Check Point のアドレスグループ変換が終了しました",
				"output", "checkpoint_address_group.conf")

			lines, err = conv_checkpoint.ConvertServices(app.Services, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_service.conf", lines, app)
			slog.Info("Check Point のサービス変換が終了しました",
				"output", "checkpoint_service.conf")

			lines, err = conv_checkpoint.ConvertServiceGroups(app.ServiceGroups, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_service_group.conf", lines, app)
			slog.Info("Check Point のサービスグループ変換が終了しました",
				"output", "checkpoint_service_group.conf")

			lines, err = conv_checkpoint.ConvertPolicies(app.Policies, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_policy.conf", lines, app)
			slog.Info("Check Point のポリシー変換が終了しました",
				"output", "checkpoint_policy.conf")

			lines, err = conv_checkpoint.ConvertNATPolicies(app.NATRules, ctx)
			if err != nil {
				slog.Error("convert error:", "err", err)
			}
			writeMgmtLines("checkpoint_nat.conf", lines, app)
			slog.Info("Check Point の NAT 変換が終了しました",
				"output", "checkpoint_nat.conf")

		default:
			slog.Error("unsupported output", "to", app.To)
		}
	case "cp":
		checkpoint.ParseCheckPoint(&app)
		slog.Info("Check Point の解析が終了しました", "output", "checkpoint.xlsx")
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
