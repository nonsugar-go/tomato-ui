package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nonsugar-go/tomato-ui/internal/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/cpdump"
)

type Config struct {
	Command  string
	Server   string
	User     string
	Password string
	Limit    int
	Output   string
}

func main() {
	handler := log.New(os.Stderr)
	handler.SetLevel(log.DebugLevel)
	handler.SetReportTimestamp(true)
	slog.SetDefault(slog.New(handler))

	cfg := parseArgs()

	client := &checkpoint.Client{
		Server: cfg.Server,
		HTTP:   checkpoint.InsecureClient(),
	}

	if err := client.Login(cfg.User, cfg.Password); err != nil {
		slog.Error("failed to login", "error", err)
		os.Exit(1)
	}

	slog.Info("login ok", "SID", client.SID)

	var data any
	var err error

	switch cfg.Command {
	case "hosts":
		data, err = cpdump.FetchHosts(client, cfg.Limit)
		if err != nil {
			slog.Error("failed to fetch hosts", "error", err)
			os.Exit(1)
		}
		slog.Info("fetched hosts", "count", len(data.([]checkpoint.Host)))
	default:
		printUsage()
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if cfg.Output != "" {
		f, err := os.Create(cfg.Output)
		if err != nil {
			slog.Error("failed to create output file", "error", err)
			os.Exit(1)
		}
		defer f.Close()
		enc = json.NewEncoder(f)
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(data); err != nil {
		slog.Error("failed to encode data", "error", err)
		os.Exit(1)
	}
}

func parseArgs() Config {
	var cfg Config

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg.Command = os.Args[1]

	fs := flag.NewFlagSet("cp-dump "+cfg.Command, flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	fs.StringVar(&cfg.Server, "server", "", "management server")
	fs.StringVar(&cfg.User, "user", "admin", "username")
	fs.StringVar(&cfg.Password, "password", "", "password")
	fs.IntVar(&cfg.Limit, "limit", 0, "number of objects to fetch (0-500, 0 = unlimited)")
	fs.StringVar(&cfg.Output, "o", "", "output file (default stdout)")

	fs.Usage = func() {
		printUsage()
		fmt.Println()
		fmt.Println("options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(os.Args[2:])
	if err != nil {
		os.Exit(1)
	}

	return cfg
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  cp-dump hosts  [options]")
	fmt.Println("  cp-dump groups [options]")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Println("  host   show-hosts from Check Point")
	fmt.Println("  group  show-groups from Check Point")
	fmt.Println()
	fmt.Println("example:")
	fmt.Println("  cp-dump host -server 192.168.1.100 -user admin -password xxx -o hosts.json")
	fmt.Println("  cp-dump group -server 192.168.1.100 -user admin -password xxx -o groups.json")
	fmt.Println()
	fmt.Println("run:")
	fmt.Println("  cp-dump hosts -h")
	fmt.Println()
	fmt.Println("jq example:")
	fmt.Println(`jq -r '.[]|[.name,.["ipv4-address"],.["ipv6-address"],.comments,(.tags//[]|join(";")),.color]|@csv' hosts.json`)
}
