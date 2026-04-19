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
	Color    string   `json:"color,omitempty"`
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
	Color    string   `json:"color,omitempty"`
}

func convertToHost(ah apiHost) Host {
	tags := make([]string, len(ah.Tags))
	for i, t := range ah.Tags {
		tags[i] = t.Name
	}
	return Host{
		Name:     ah.Name,
		IPv4:     ah.IPv4,
		IPv6:     ah.IPv6,
		Comments: ah.Comments,
		Tags:     tags,
		Color:    ah.Color,
	}
}

// ShowHosts retrieves a list of hosts from the Check Point management server.
// Ref: https://sc1.checkpoint.com/documents/latest/APIs/?#web/show-hosts
func (c *Client) ShowHosts(limit int) ([]Host, error) {
	const pageSize = 500

	var (
		offset int
		all    []Host
	)

	for {
		currentLimit := limit

		if limit == 0 {
			currentLimit = pageSize
		}

		body := map[string]any{
			"limit":         currentLimit,
			"offset":        offset,
			"details-level": "full",
		}

		resp, err := c.Post("show-hosts", body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}

		var result struct {
			Objects []apiHost `json:"objects"`
			Total   int       `json:"total"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		for _, ah := range result.Objects {
			all = append(all, convertToHost(ah))
		}

		if limit > 0 {
			break
		}

		offset += len(result.Objects)
		if len(result.Objects) == 0 {
			break
		}
	}

	return all, nil
}
