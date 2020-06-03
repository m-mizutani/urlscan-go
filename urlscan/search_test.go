package urlscan_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/urlscan-go/urlscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)

	resp, err := client.Search(context.Background(), urlscan.SearchArguments{
		Query: urlscan.String("ip:163.43.24.70"),
	})

	require.NoError(t, err)
	assert.NotEqual(t, 0, len(resp.Results))
	assert.NotEqual(t, "", resp.Results[0].ID)
}

func TestSearchSize(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)

	resp, err := client.Search(context.Background(), urlscan.SearchArguments{
		Query: urlscan.String("ip:163.43.24.70"),
		Size:  urlscan.Uint64(1),
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(resp.Results))
}

func TestSearchOffset(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)

	resp1, err := client.Search(context.Background(), urlscan.SearchArguments{
		Query:  urlscan.String("ip:163.43.24.70"),
		Size:   urlscan.Uint64(1),
		Offset: urlscan.Uint64(0),
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(resp1.Results))

	resp2, err := client.Search(context.Background(), urlscan.SearchArguments{
		Query:  urlscan.String("ip:163.43.24.70"),
		Size:   urlscan.Uint64(1),
		Offset: urlscan.Uint64(1),
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(resp2.Results))

	assert.NotEqual(t, resp1.Results[0].ID, resp2.Results[0].ID)
}
