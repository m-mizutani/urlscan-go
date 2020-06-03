package urlscan

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// SubmitArguments is input argument of Submit()
type SubmitArguments struct {
	// URL is Required.
	URL string `json:"url"`
	// CustomAgent is optional. You can set own UserAgent as you like it.
	CustomAgent *string `json:"customagent"`
	// Referer is optional. You can set any Referer.
	Referer *string `json:"referer"`
	// Public is optional. Default is "off" that means "private". You need to set it as "on" if you want to make the result public.
	Public *string `json:"public"`
}

type submitResponse struct {
	Visibility string `json:"visibility"`
	URL        string `json:"url"`
	Message    string `json:"message"`
	UUID       string `json:"uuid"`
	Result     string `json:"result"`
	API        string `json:"api"`
	Options    map[string]interface{}
}

// Submit sends a request of sandbox execution for specified URL.
func (x *Client) Submit(ctx context.Context, args SubmitArguments) (Task, error) {
	task := Task{
		client: x,
	}

	var result submitResponse
	code, err := x.post(ctx, "scan", args, &result)
	if err != nil {
		return task, err
	}
	if code != 200 {
		return task, errors.Errorf("Unexpected status code: %d", code)
	}

	task.url = result.API
	task.uuid = result.UUID
	return task, nil
}

// Task is returned by Submit() and you can fetch a result of the submitted scan from the Task.
type Task struct {
	client *Client
	uuid   string
	url    string
	Result ScanResult
}

func getExpWaitTime(count int) time.Duration {
	d := (count * count * 100) + 1000
	if d > 20*1000 {
		d = 20 * 1000
	}

	return time.Millisecond * time.Duration(d)
}

// Wait tries to retrieve a result. If scan is not still completed, it retries up to 30 times
func (x *Task) WaitForReport(ctx context.Context) error {
	maxRetry := 30
	for i := 0; i < maxRetry; i++ {
		code, err := x.client.get(ctx, fmt.Sprintf("result/%s", x.uuid), nil, &x.Result)
		if err != nil {
			return errors.Wrap(err, "Fail to get result query")
		}

		switch code {
		case 200:
			return nil
		case 400:
			return errors.New("status: 400")
		}

		time.Sleep(getExpWaitTime(i))
	}

	return errors.Errorf("Timeout of task id: %s", x.uuid)
}

// -------------------------------
// Structures of scan result
// -------------------------------

// ScanResult of a root data structure of scan result.
type ScanResult struct {
	Data  ScanData  `json:"data"`
	Lists ScanLists `json:"lists"`
	Meta  ScanMeta  `json:"meta"`
	Page  ScanPage  `json:"page"`
	Stats ScanStats `json:"stats"`
	Task  ScanTask  `json:"task"`
}

// ScanGeo presents GeoLocation information
type ScanGeo struct {
	City        string      `json:"city"`
	Country     string      `json:"country"`
	CountryName string      `json:"country_name"`
	LL          []float64   `json:"ll"`
	Metro       int64       `json:"metro"`
	Range       interface{} `json:"range"`
	Region      string      `json:"region"`
	Zip         int64       `json:"zip"`
}

// ScanData presents main result of the scan
type ScanData struct {
	Console []interface{} `json:"console"`

	Cookies []struct {
		Domain   string  `json:"domain"`
		Expires  float64 `json:"expires"`
		HTTPOnly bool    `json:"httpOnly"`
		Name     string  `json:"name"`
		Path     string  `json:"path"`
		Secure   bool    `json:"secure"`
		Session  bool    `json:"session"`
		Size     int64   `json:"size"`
		Value    string  `json:"value"`
	} `json:"cookies"`

	Globals []struct {
		Prop string `json:"prop"`
		Type string `json:"type"`
	} `json:"globals"`
	Links []struct {
		Href string `json:"href"`
		Text string `json:"text"`
	} `json:"links"`

	Requests []struct {
		InitiatorInfo struct {
			Host string `json:"host"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"initiatorInfo"`

		Request struct {
			DocumentURL    string `json:"documentURL"`
			FrameID        string `json:"frameId"`
			HasUserGesture bool   `json:"hasUserGesture"`
			Initiator      struct {
				LineNumber int64 `json:"lineNumber"`
				Stack      struct {
					CallFrames []struct {
						ColumnNumber int64  `json:"columnNumber"`
						FunctionName string `json:"functionName"`
						LineNumber   int64  `json:"lineNumber"`
						ScriptID     string `json:"scriptId"`
						URL          string `json:"url"`
					} `json:"callFrames"`
				} `json:"stack"`
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"initiator"`
			LoaderID string `json:"loaderId"`
			Request  struct {
				HasPostData bool `json:"hasPostData"`
				Headers     struct {
					ContentType             string `json:"Content-Type"`
					Origin                  string `json:"Origin"`
					Referer                 string `json:"Referer"`
					UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
					UserAgent               string `json:"User-Agent"`
				} `json:"headers"`
				InitialPriority  string `json:"initialPriority"`
				Method           string `json:"method"`
				MixedContentType string `json:"mixedContentType"`
				PostData         string `json:"postData"`
				ReferrerPolicy   string `json:"referrerPolicy"`
				URL              string `json:"url"`
			} `json:"request"`
			RequestID string  `json:"requestId"`
			Timestamp float64 `json:"timestamp"`
			Type      string  `json:"type"`
			WallTime  float64 `json:"wallTime"`
		} `json:"request"`

		Response struct {
			Abp struct {
				Source string `json:"source"`
				Type   string `json:"type"`
				URL    string `json:"url"`
			} `json:"abp"`
			Asn struct {
				Asn         string `json:"asn"`
				Country     string `json:"country"`
				Date        string `json:"date"`
				Description string `json:"description"`
				IP          string `json:"ip"`
				Name        string `json:"name"`
				Registrar   string `json:"registrar"`
				Route       string `json:"route"`
			} `json:"asn"`
			DataLength        int64         `json:"dataLength"`
			EncodedDataLength int64         `json:"encodedDataLength"`
			Geoip             ScanGeo       `json:"geoip"`
			Hash              string        `json:"hash"`
			Hashmatches       []interface{} `json:"hashmatches"`
			Rdns              struct {
				IP  string `json:"ip"`
				Ptr string `json:"ptr"`
			} `json:"rdns"`
			RequestID string `json:"requestId"`

			Response struct {
				EncodedDataLength int64             `json:"encodedDataLength"`
				Headers           map[string]string `json:"headers"`
				MimeType          string            `json:"mimeType"`
				Protocol          string            `json:"protocol"`
				RemoteIPAddress   string            `json:"remoteIPAddress"`
				RemotePort        int64             `json:"remotePort"`
				RequestHeaders    map[string]string `json:"requestHeaders"`
				SecurityDetails   struct {
					CertificateID                     int64         `json:"certificateId"`
					CertificateTransparencyCompliance string        `json:"certificateTransparencyCompliance"`
					Cipher                            string        `json:"cipher"`
					Issuer                            string        `json:"issuer"`
					KeyExchange                       string        `json:"keyExchange"`
					KeyExchangeGroup                  string        `json:"keyExchangeGroup"`
					Protocol                          string        `json:"protocol"`
					SanList                           []string      `json:"sanList"`
					SignedCertificateTimestampList    []interface{} `json:"signedCertificateTimestampList"`
					SubjectName                       string        `json:"subjectName"`
					ValidFrom                         int64         `json:"validFrom"`
					ValidTo                           int64         `json:"validTo"`
				} `json:"securityDetails"`
				SecurityHeaders []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"securityHeaders"`
				SecurityState string `json:"securityState"`
				Status        int64  `json:"status"`
				StatusText    string `json:"statusText"`
				Timing        struct {
					ConnectEnd        float64 `json:"connectEnd"`
					ConnectStart      float64 `json:"connectStart"`
					DNSEnd            float64 `json:"dnsEnd"`
					DNSStart          float64 `json:"dnsStart"`
					ProxyEnd          int64   `json:"proxyEnd"`
					ProxyStart        int64   `json:"proxyStart"`
					PushEnd           int64   `json:"pushEnd"`
					PushStart         int64   `json:"pushStart"`
					ReceiveHeadersEnd float64 `json:"receiveHeadersEnd"`
					RequestTime       float64 `json:"requestTime"`
					SendEnd           float64 `json:"sendEnd"`
					SendStart         float64 `json:"sendStart"`
					SslEnd            float64 `json:"sslEnd"`
					SslStart          float64 `json:"sslStart"`
					WorkerReady       int64   `json:"workerReady"`
					WorkerStart       int64   `json:"workerStart"`
				} `json:"timing"`
				URL string `json:"url"`
			} `json:"response"`
			Size int64  `json:"size"`
			Type string `json:"type"`
		} `json:"response"`
	} `json:"requests"`
	Timing struct {
		BeginNavigation      string `json:"beginNavigation"`
		DomContentEventFired string `json:"domContentEventFired"`
		FrameNavigated       string `json:"frameNavigated"`
		FrameStartedLoading  string `json:"frameStartedLoading"`
		FrameStoppedLoading  string `json:"frameStoppedLoading"`
		LoadEventFired       string `json:"loadEventFired"`
	} `json:"timing"`
}

// ScanLists shows lists
type ScanLists struct {
	Asns         []string `json:"asns"`
	Certificates []struct {
		Issuer      string `json:"issuer"`
		SubjectName string `json:"subjectName"`
		ValidFrom   int64  `json:"validFrom"`
		ValidTo     int64  `json:"validTo"`
	} `json:"certificates"`
	Countries   []string      `json:"countries"`
	Domains     []string      `json:"domains"`
	Hashes      []interface{} `json:"hashes"`
	Ips         []string      `json:"ips"`
	LinkDomains []string      `json:"linkDomains"`
	Servers     []string      `json:"servers"`
	Urls        []string      `json:"urls"`
}

// ScanMeta presents scan meta data
type ScanMeta struct {
	Processors struct {
		Abp struct {
			Data []struct {
				Source string `json:"source"`
				Type   string `json:"type"`
				URL    string `json:"url"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"abp"`
		Asn struct {
			Data []struct {
				Asn         string `json:"asn"`
				Country     string `json:"country"`
				Date        string `json:"date"`
				Description string `json:"description"`
				IP          string `json:"ip"`
				Name        string `json:"name"`
				Registrar   string `json:"registrar"`
				Route       string `json:"route"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"asn"`
		Cdnjs struct {
			Data  []interface{} `json:"data"`
			State string        `json:"state"`
		} `json:"cdnjs"`
		Done struct {
			Data struct {
				State string `json:"state"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"done"`
		Download struct {
			Data  []interface{} `json:"data"`
			State string        `json:"state"`
		} `json:"download"`
		Geoip struct {
			Data []struct {
				Geoip ScanGeo `json:"geoip"`
				IP    string  `json:"ip"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"geoip"`
		Gsb struct {
			Data  struct{} `json:"data"`
			State string   `json:"state"`
		} `json:"gsb"`
		Rdns struct {
			Data []struct {
				IP  string `json:"ip"`
				Ptr string `json:"ptr"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"rdns"`
		Wappa struct {
			Data []struct {
				App        string `json:"app"`
				Categories []struct {
					Name     string `json:"name"`
					Priority string `json:"priority"`
				} `json:"categories"`
				Confidence []struct {
					Confidence interface{} `json:"confidence"`
					Pattern    string      `json:"pattern"`
				} `json:"confidence"`
				ConfidenceTotal int64  `json:"confidenceTotal"`
				Icon            string `json:"icon"`
				Website         string `json:"website"`
			} `json:"data"`
			State string `json:"state"`
		} `json:"wappa"`
	} `json:"processors"`
}

// ScanPage shows page information
type ScanPage struct {
	Asn     string `json:"asn"`
	Asnname string `json:"asnname"`
	City    string `json:"city"`
	Country string `json:"country"`
	Domain  string `json:"domain"`
	IP      string `json:"ip"`
	Ptr     string `json:"ptr"`
	Server  string `json:"server"`
	URL     string `json:"url"`
}

// ScanTask presents submitted task of scan
type ScanTask struct {
	DomURL  string `json:"domURL"`
	Method  string `json:"method"`
	Options struct {
		Useragent string `json:"useragent"`
	} `json:"options"`
	ReportURL     string `json:"reportURL"`
	ScreenshotURL string `json:"screenshotURL"`
	Source        string `json:"source"`
	Time          string `json:"time"`
	URL           string `json:"url"`
	UserAgent     string `json:"userAgent"`
	UUID          string `json:"uuid"`
	Visibility    string `json:"visibility"`
}

// ScanStatsDetail is a detail of scan
type ScanStatsDetail struct {
	Compression   string      `json:"compression"`
	Count         int64       `json:"count"`
	Countries     []string    `json:"countries"`
	Domain        string      `json:"domain"`
	EncodedSize   int64       `json:"encodedSize"`
	Index         int64       `json:"index"`
	Initiators    []string    `json:"initiators"`
	Ips           []string    `json:"ips"`
	Latency       int64       `json:"latency"`
	Percentage    interface{} `json:"percentage"`
	Protocol      string      `json:"protocol"`
	Protocols     interface{} `json:"protocols"`
	Redirects     int64       `json:"redirects"`
	RegDomain     string      `json:"regDomain"`
	SecurityState interface{} `json:"securityState"`
	Server        string      `json:"server"`
	Size          int64       `json:"size"`
	Type          string      `json:"type"`
	SubDomains    []struct {
		Domain string `json:"domain"`
		Failed bool   `json:"failed"`
	} `json:"subDomains"`
}

// ScanStats is a like summary of scan.
type ScanStats struct {
	IPv6Percentage int64             `json:"IPv6Percentage"`
	AdBlocked      int64             `json:"adBlocked"`
	DomainStats    []ScanStatsDetail `json:"domainStats"`
	IPStats        []struct {
		Asn struct {
			Asn         string `json:"asn"`
			Country     string `json:"country"`
			Date        string `json:"date"`
			Description string `json:"description"`
			IP          string `json:"ip"`
			Name        string `json:"name"`
			Registrar   string `json:"registrar"`
			Route       string `json:"route"`
		} `json:"asn"`
		Count       interface{} `json:"count"`
		Countries   []string    `json:"countries"`
		DNS         struct{}    `json:"dns"`
		Domains     []string    `json:"domains"`
		EncodedSize int64       `json:"encodedSize"`
		Geoip       ScanGeo     `json:"geoip"`
		Index       int64       `json:"index"`
		IP          string      `json:"ip"`
		Ipv6        bool        `json:"ipv6"`
		Rdns        struct {
			IP  string `json:"ip"`
			Ptr string `json:"ptr"`
		} `json:"rdns"`
		Redirects int64 `json:"redirects"`
		Requests  int64 `json:"requests"`
		Size      int64 `json:"size"`
	} `json:"ipStats"`
	Malicious        int64             `json:"malicious"`
	ProtocolStats    []ScanStatsDetail `json:"protocolStats"`
	RegDomainStats   []ScanStatsDetail `json:"regDomainStats"`
	ResourceStats    []ScanStatsDetail `json:"resourceStats"`
	SecurePercentage int64             `json:"securePercentage"`
	SecureRequests   int64             `json:"secureRequests"`
	ServerStats      []ScanStatsDetail `json:"serverStats"`
	TLSStats         []ScanStatsDetail `json:"tlsStats"`
	TotalLinks       int64             `json:"totalLinks"`
	UniqCountries    int64             `json:"uniqCountries"`
}
