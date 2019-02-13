package urlscan_test

import (
	"log"
	"os"
	"testing"

	"github.com/m-mizutani/urlscan-go/urlscan"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type config struct {
	ApiKey string `json:"api_key"`
}

var cfg config

func init() {
	cfg.ApiKey = os.Getenv("URLSCAN_API_KEY")
	if cfg.ApiKey == "" {
		log.Fatal("no API KEY, environment variable URLSCAN_API_KEY is required.")
	}

	urlscan.Logger = logrus.New()
	urlscan.Logger.SetLevel(logrus.InfoLevel)
}

func TestSubmitScan(t *testing.T) {
	client := urlscan.NewClient(cfg.ApiKey)
	task, err := client.Submit(urlscan.SubmitArguments{
		URL: "https://cookpad.com",
	})

	require.NoError(t, err)
	err = task.Wait()
	require.NoError(t, err)
}
