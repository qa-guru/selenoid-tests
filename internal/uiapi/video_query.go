package uiapi

import (
	"fmt"
	"net/url"
)

func buildVideoQuery(limit, offset int, q string) url.Values {
	query := url.Values{}
	query.Set("json", "")
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))
	if q != "" {
		query.Set("q", q)
	}
	return query
}
