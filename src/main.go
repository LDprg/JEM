package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

const version = "0.1.0"

var osName = runtime.GOOS
var arch = getArch()

func getArch() string {
	var arch = runtime.GOARCH

	switch arch {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	}

	return arch
}

func list() {
	fmt.Println("list")
}

func use() {
	fmt.Println("use")
}

func install(version string) {
	fmt.Println("Installing JDK version: " + version)

	fmt.Println("Downloading...")
	//fmt.Println("Arch: " + arch)
	//fmt.Println("OS: " + osName)

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

	fileName := "data/installed/jdk-" + version + ".zip"

	file, err := os.Create(fileName)
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

	fmt.Println("Extracting...")

	unzip(fileName, "data/current/jdk-"+version)
}

func uninstall() {
	fmt.Println("uninstall")
}

func help() {
	fmt.Println("Runnning version: " + version)
	fmt.Println("Usage: lem [command]")
	fmt.Println("Commands:")
	fmt.Println("  list")
	fmt.Println("  use [version]")
	fmt.Println("  install [version]")
	fmt.Println("  uninstall [version]")
}
func selectCommand() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			list()
		case "use":
			use()
		case "install":
			if len(os.Args) > 2 {
				install(os.Args[2])
			} else {
				help()
			}
		case "uninstall":
			uninstall()
		default:
			help()
		}
	} else {
		help()
	}
}

func main() {
	os.Args = append(os.Args, "install")
	os.Args = append(os.Args, "8")

	temp, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Mkdir("data", os.ModePerm)
	_ = os.Mkdir("data/current", os.ModePerm)
	_ = os.Mkdir("data/installed", os.ModePerm)

	selectCommand()

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(temp)
}
