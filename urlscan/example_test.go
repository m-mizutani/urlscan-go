package urlscan_test

import (
	"fmt"
	"log"

	"github.com/m-mizutani/urlscan-go/urlscan"
)

func ExampleClient_Submit() {
	client := urlscan.NewClient("YOUR-API-KEY")
	task, err := client.Submit(urlscan.SubmitArguments{URL: "https://golang.org"})
	if err != nil {
		log.Fatal(err)
	}

	err = task.Wait()
	if err != nil {
		log.Fatal(err)
	}

	for _, cookie := range task.Result.Data.Cookies {
		fmt.Printf("Cookie: %s = %s\n", cookie.Name, cookie.Value)
	}
}

func ExampleClient_Search() {
	client := urlscan.NewClient("YOUR-API-KEY")

	resp, err := client.Search(urlscan.SearchArguments{
		Query:  urlscan.String("ip:1.2.3.x"),
		Size:   urlscan.Uint64(1),
		Offset: urlscan.Uint64(0),
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("Related URL: %s\n", result.Page.URL)
	}
}
