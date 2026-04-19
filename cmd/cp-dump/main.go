package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nonsugar-go/tomato-ui/internal/checkpoint"
	"github.com/nonsugar-go/tomato-ui/internal/cpdump"
)

type Config struct {
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
	data, err := cpdump.FetchHosts(client, cfg.Limit)
	if err != nil {
		slog.Error("failed to fetch hosts", "error", err)
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

	flag.StringVar(&cfg.Server, "server", "", "management server")
	flag.StringVar(&cfg.User, "user", "admin", "username")
	flag.StringVar(&cfg.Password, "password", "", "password")
	flag.IntVar(&cfg.Limit, "limit", 500, "fetch limit")
	flag.StringVar(&cfg.Output, "o", "", "output file (default stdout)")

	flag.Parse()

	if cfg.Server == "" {
		slog.Error("server is required")

		slog.Info(`example: cp-dump -server 192.168.1.41 -user secadmin -password Lab@12345 -o hosts.json`)
		slog.Info(`example: jq -r '.[] | [.name, .["ipv4-address"], .["ipv6-address"], .comments, (.tags // [] | join(";")) ] | @csv' hosts.json`)
		os.Exit(1)
	}

	return cfg
}
