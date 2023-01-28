package main

import (
	"fmt"
	"github.com/YakirOren/sdarot"
	"log"
	"os"
)

func main() {
	client, err := sdarot.New(sdarot.Config{
		Username: "username",
		Password: "password",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Getting video")
	video, err := client.GetVideo(sdarot.VideoRequest{
		SeriesID: "19",
		Season:   "1",
		Episode:  "1",
	})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(video.ID + ".mp4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading", file.Name())

	if err := client.Download(video, file); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
