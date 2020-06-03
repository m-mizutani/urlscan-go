package urlscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logger is a logrus logger. You can replace the logger with yours or change setting if you need.
var Logger = logrus.New()

func init() {
	Logger.SetLevel(logrus.PanicLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	Logger.SetReportCaller(true)
}

// String converts string variable and literal to pointer
func String(s string) *string {
	return &s
}

// Uint64 converts uint64 variable and literal to pointer
func Uint64(u uint64) *uint64 {
	return &u
}

// Client is main structure of the library, a requester to urlscan.io.
type Client struct {
	apiKey  string
	BaseURL string
}

// NewClient is a constructor of Client
func NewClient(apiKey string) Client {
	client := Client{
		apiKey:  apiKey,
		BaseURL: "https://urlscan.io/api/v1",
	}

	return client
}

func (x Client) post(ctx context.Context, apiName string, input interface{}, output interface{}) (int, error) {
	rawData, err := json.Marshal(input)
	if err != nil {
		return 0, errors.Wrap(err, "Fail to marshal urlscan.io submit argument")
	}

	uri := fmt.Sprintf("%s/%s/", x.BaseURL, apiName)
	Logger.WithFields(logrus.Fields{
		"uri":  uri,
		"body": string(rawData),
	}).Debug("Generated Query")

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewReader(rawData))
	if err != nil {
		return 0, errors.Wrap(err, "Fail to create urlscan.io scan POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("API-Key", x.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return -1, errors.Wrap(err, "Fail to send urlscan.io POST request")
		}
		return resp.StatusCode, errors.Wrap(err, "Fail to send urlscan.io POST request")
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, errors.Wrap(err, "Fail to read urlscan.io POST result")
	}
	if resp.StatusCode != 200 {
		Logger.WithFields(logrus.Fields{
			"body": string(buf),
			"code": resp.StatusCode,
		}).Warn("Unexpected status code")
	}

	err = json.Unmarshal(buf, &output)
	if err != nil {
		return resp.StatusCode, errors.Wrap(err, "Fail to unmarshal urlscan.io POST result")
	}

	return resp.StatusCode, nil
}

func (x Client) get(ctx context.Context, apiName string, values url.Values, output interface{}) (int, error) {
	var qs string
	if values != nil {
		qs = "?" + values.Encode()
	}

	uri := fmt.Sprintf("%s/%s/%s", x.BaseURL, apiName, qs)
	Logger.WithField("uri", uri).Info("Generated Query")

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return 0, errors.Wrap(err, "Fail to create urlscan.io get request")
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return -1, errors.Wrap(err, "Fail to send urlscan.io POST request")
		}
		return resp.StatusCode, errors.Wrap(err, "Fail to send urlscan.io get request")
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, errors.Wrap(err, "Fail to read urlscan.io get result")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		Logger.WithFields(logrus.Fields{
			"body": string(buf),
			"code": resp.StatusCode,
		}).Warn("Unexpected status code")
	}

	err = json.Unmarshal(buf, &output)
	if err != nil {
		return resp.StatusCode, errors.Wrap(err, "Fail to unmarshal urlscan.io get result")
	}

	return resp.StatusCode, nil
}
