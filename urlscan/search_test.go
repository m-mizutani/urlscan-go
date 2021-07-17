package urlscan_test

import (
	"testing"

	"github.com/m-mizutani/urlscan-go/urlscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)

	resp, err := client.Search(urlscan.SearchArguments{
		Query: urlscan.String("ip:163.43.24.70"),
	})

	require.NoError(t, err)
	assert.NotEqual(t, 0, len(resp.Results))
	assert.NotEqual(t, "", resp.Results[0].ID)
}

func TestSearchSize(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)

	resp, err := client.Search(urlscan.SearchArguments{
		Query: urlscan.String("ip:163.43.24.70"),
		Size:  urlscan.Uint64(1),
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(resp.Results))
}
