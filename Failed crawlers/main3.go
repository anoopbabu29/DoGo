package main3

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "DoGo-80cbcac25e42.json")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "dogo-236814")

	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}
	urls = make(map[string]bool)

	words = make(map[string]int)

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)

	fmt.Println(client)
	if err != nil {
		fmt.Println("client")
	}

	data, err := ioutil.ReadFile("data-download-pub78.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(data))
}
