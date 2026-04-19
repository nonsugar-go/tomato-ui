package cpdump

import "github.com/nonsugar-go/tomato-ui/internal/checkpoint"

func FetchHosts(client *checkpoint.Client, limit int) ([]checkpoint.Host, error) {
	return client.ShowHosts(limit)
}
