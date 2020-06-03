package urlscan

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// SearchArguments is input data structure of Search()
type SearchArguments struct {
	// Optional. urlscan.io search query.
	// See Help & Example of https://urlscan.io/search/ for more detail
	Query *string `json:"query"`
	// Optional. Page size
	Size *uint64 `json:"size"`
	// Optional. Offset of the result
	Offset *uint64 `json:"offset"`
	// Optional. specificied via $sort_field:$sort_order. Default: _score
	Sort *string `json:"sort"`
}

// SearchResult represents a single search result from the API
type SearchResult struct {
	ID   string `json:"_id"`
	Page struct {
		Asn     string `json:"asn"`
		Asnname string `json:"asnname"`
		City    string `json:"city"`
		Country string `json:"country"`
		Domain  string `json:"domain"`
		IP      string `json:"ip"`
		Ptr     string `json:"ptr"`
		Server  string `json:"server"`
		URL     string `json:"url"`
	} `json:"page"`
	Result string `json:"result"`
	Stats  struct {
		ConsoleMsgs       int64 `json:"consoleMsgs"`
		DataLength        int64 `json:"dataLength"`
		EncodedDataLength int64 `json:"encodedDataLength"`
		Requests          int64 `json:"requests"`
		UniqIPs           int64 `json:"uniqIPs"`
	} `json:"stats"`
	Task struct {
		Method     string `json:"method"`
		Source     string `json:"source"`
		Time       string `json:"time"`
		URL        string `json:"url"`
		Visibility string `json:"visibility"`
	} `json:"task"`
	UniqCountries int64 `json:"uniq_countries"`
}

// SearchResponse is returned by Search() and including existing scan results.
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int64          `json:"total"`
}

// Search sends query to search existing scan results with query
func (x *Client) Search(ctx context.Context, args SearchArguments) (SearchResponse, error) {
	var result SearchResponse
	values := make(url.Values)

	if args.Query != nil {
		values.Add("q", *args.Query)
	}
	if args.Size != nil {
		values.Add("size", fmt.Sprintf("%d", *args.Size))
	}
	if args.Offset != nil {
		values.Add("offset", fmt.Sprintf("%d", *args.Offset))
	}
	if args.Sort != nil {
		values.Add("sort", *args.Sort)
	}

	code, err := x.get(ctx, "search", values, &result)
	if err != nil {
		return result, err
	}
	if code != 200 {
		return result, errors.Errorf("Unexpected status code: %d", code)
	}

	return result, err
}
