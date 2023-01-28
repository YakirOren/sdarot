package main

import (
	"fmt"
	"github.com/YakirOren/sdarot"
	"log"
)

func main() {
	client, err := sdarot.New(sdarot.Config{
		Username: "username",
		Password: "password",
	})
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.Search("the b")
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result.EnglishName, "Series ID:", result.SeriesID)
	}
}
