package main

import (
	"fmt"
	"log"
	"os"

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

	fmt.Println("Getting video")
	video, err := client.GetVideo(sdarot.VideoRequest{
		SeriesID: 19,
		Season:   1,
		Episode:  1,
	})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("%d.mp4", video.ID))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading", file.Name())

	if err := client.Download(video, file); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
