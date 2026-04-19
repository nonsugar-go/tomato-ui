package checkpoint

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Name     string   `json:"name"`
	IPv4     string   `json:"ipv4-address,omitempty"`
	IPv6     string   `json:"ipv6-address,omitempty"`
	Comments string   `json:"comments,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type apiTag struct {
	Name string `json:"name"`
}

type apiHost struct {
	Name     string   `json:"name"`
	IPv4     string   `json:"ipv4-address,omitempty"`
	IPv6     string   `json:"ipv6-address,omitempty"`
	Comments string   `json:"comments,omitempty"`
	Tags     []apiTag `json:"tags,omitempty"`
}

// ShowHosts retrieves a list of hosts from the Check Point management server.
// Ref: https://sc1.checkpoint.com/documents/latest/APIs/?#web/show-hosts
func (c *Client) ShowHosts(limit int) ([]Host, error) {
	body := map[string]any{
		"limit":         limit,
		"details-level": "full",
	}

	resp, err := c.Post("show-hosts", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result struct {
		Objects []apiHost `json:"objects"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var hosts []Host
	for _, h := range result.Objects {
		var tags []string
		for _, t := range h.Tags {
			tags = append(tags, t.Name)
		}

		hosts = append(hosts, Host{
			Name:     h.Name,
			IPv4:     h.IPv4,
			IPv6:     h.IPv6,
			Comments: h.Comments,
			Tags:     tags,
		})
	}

	return hosts, nil
}
