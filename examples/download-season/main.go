package main

import (
	"errors"
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

	series, err := client.GetSeriesByID(6746)
	if err != nil {
		log.Fatal(err)
	}

	for _, ep := range series.GetEpisodes(1) {
		fmt.Println("Getting video: ", ep.Episode)

		video, err := client.GetVideo(ep)

		if errors.Is(err, sdarot.ErrServerOverLoad) {
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		download(series.EnglishName, video, client)

	}
}

func download(folderName string, video *sdarot.Video, client *sdarot.Client) {
	path := fmt.Sprintf("%s/%d", folderName, video.Metadata.Season)
	err := os.MkdirAll(path, 0o777)
	if err != nil {
		log.Println(err)
	}
	file, err := os.Create(fmt.Sprintf("%s/%d.mp4", path, video.Metadata.Episode))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading", file.Name())

	if err := client.Download(video, file); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
