package main

import (
	"encoding/json"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

type availableJAVA struct {
	AvailableReleases []int `json:"available_releases"`
}

func available() []int {
	url := "https://api.adoptium.net/v3/info/available_releases"
	get, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(get.Body)

	out, err := io.ReadAll(get.Body)
	if err != nil {
		log.Fatal(err)
	}

	var av availableJAVA

	err = json.Unmarshal(out, &av)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(av.AvailableReleases[:], func(i, j int) bool {
		return av.AvailableReleases[:][i] > av.AvailableReleases[:][j]
	})

	return av.AvailableReleases
}

func download(version string, osName string, arch string, dest string) {
	url := "https://api.adoptium.net/v3/binary/latest/" + version + "/ga/" + osName + "/" + arch + "/jdk/hotspot/normal/eclipse?project=jdk"

	get, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(get.Body)

	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}

	filesize := get.ContentLength

	bar := pb.Full.Start64(filesize)
	barReader := bar.NewProxyReader(get.Body)

	_, err = io.Copy(file, barReader)
	if err != nil {
		return
	}

	bar.Finish()
}
