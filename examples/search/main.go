package main

import (
	"fmt"
	"log"

	"github.com/YakirOren/sdarot"
)

func main() {
	client, err := sdarot.New(sdarot.Config{
		Username: "user",
		Password: "Password1",
	})
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.Search("the b")
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result.EnglishName, "series ID:", result.ID)
	}
}
