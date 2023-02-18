# Sdarot

[![GoDoc](https://godoc.org/github.com/yakiroren/sdarot?status.svg)](https://godoc.org/github.com/yakiroren/sdarot)

Wrapper for sdarot.tv .

* [examples](https://github.com/YakirOren/sdarot/tree/main/examples)

## Installation

```
go get github.com/yakiroren/sdarot
```

## Example Usage

```go
client, _ := sdarot.New(sdarot.Config{
    Username: "user",
    Password: "Password1",
})

// Get episode
video, _ := client.GetVideo(sdarot.VideoRequest{
    SeriesID: 19,
    Season:   1,
    Episode:  1,
})

// save video to file
file, _ := os.Create(fmt.Sprintf("%d.mp4", video.ID))

client.Download(video, file)
```

