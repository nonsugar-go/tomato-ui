package checkpoint

import (
	"encoding/json"
	"fmt"
)

type apiServiceTCP struct {
	Name        string `json:"name"`
	Port        string `json:"port,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Comments    string `json:"comments,omitempty"`
	MatchByType string `json:"match-by-type,omitempty"`
	Tags        []struct {
		Name string `json:"name"`
	} `json:"tags,omitempty"`
}

type ServiceTCP struct {
	Name        string
	Port        string
	Protocol    string
	Comments    string
	MatchByType string
	Tags        []string
}

// ShowServiceTCP retrieves TCP services from Check Point management server.
// Ref: https://sc1.checkpoint.com/documents/latest/APIs/?#web/show-services-tcp
func (c *Client) ShowServiceTCP(limit int) ([]ServiceTCP, error) {
	const pageSize = 500

	var (
		offset int
		all    []ServiceTCP
	)

	for {
		currentLimit := limit

		// unlimited mode
		if limit == 0 {
			currentLimit = pageSize
		}

		body := map[string]any{
			"limit":         currentLimit,
			"offset":        offset,
			"details-level": "full",
		}

		resp, err := c.Post("show-services-tcp", body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}

		var result struct {
			Objects []apiServiceTCP `json:"objects"`
			Total   int             `json:"total"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		for _, as := range result.Objects {
			all = append(all, convertToServiceTCP(as))
		}

		// finite mode
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

func convertToServiceTCP(a apiServiceTCP) ServiceTCP {
	var tags []string
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}

	return ServiceTCP{
		Name:        a.Name,
		Port:        a.Port,
		Protocol:    a.Protocol,
		Comments:    a.Comments,
		MatchByType: a.MatchByType,
		Tags:        tags,
	}
}
