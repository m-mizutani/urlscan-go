package urlscan

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type SearchArguments struct {
	Query  *string `json:"query"`
	Size   *uint64 `json:"size"`
	Offset *uint64 `json:"offset"`
	Sort   *string `json:"sort"`
}

type SearchResponse struct {
	Results []struct {
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
	} `json:"results"`
	Total int64 `json:"total"`
}

func (x *Client) Search(args SearchArguments) (SearchResponse, error) {
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

	code, err := x.get("search", values, &result)
	if err != nil {
		return result, err
	}
	if code != 200 {
		return result, errors.Errorf("Unexpected status code: %d", code)
	}

	return result, err
}
