# urlscan.io in Go 

[![Report card](https://goreportcard.com/badge/github.com/m-mizutani/urlscan-go)](https://goreportcard.com/report/github.com/m-mizutani/urlscan-go) &nbsp; [![Build Status](https://travis-ci.org/m-mizutani/urlscan-go.svg?branch=master)](https://travis-ci.org/m-mizutani/urlscan-go)&nbsp;[![GoDoc](https://godoc.org/github.com/m-mizutani/urlscan-go?status.svg)](https://godoc.org/github.com/m-mizutani/urlscan-go)

The package provides a API client of [urlscan.io](https://urlscan.io) in Go.

```go
package main

import (
  "github.com/m-mizutani/urlscan-go/urlscan"
  "fmt"
)

func main() {
    client := urlscan.NewClient("YOUR-API-KEY")
    task, err := client.Submit(urlscan.SubmitArguments{URL: "https://golang.org"})
    if err != nil {
        log.Fatal(err)
    }

    err = task.WaitForReport(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for _, cookie := range task.Result.Data.Cookies {
        fmt.Printf("Cookie: %s = %s\n", cookie.Name, cookie.Value)
    }
}
```

## Document

https://godoc.org/github.com/m-mizutani/urlscan-go/urlscan

## Test

You need to retrieve API key at first. See https://urlscan.io/about-api/#integrations for more detail.

```bash
env URLSCAN_API_KEY=12345678-your-apikey go test ./urlscan
```

## License

- Author: Masayoshi Mizutani <mizutani@sfc.wide.ad.jp>
- [The 3-Clause BSD License](./LICENSE)
