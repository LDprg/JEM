package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const version = "0.1.0"
const data = "data"

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
	fmt.Println("Current JAVA version:")

	current, err := os.ReadFile(data + "/current_version.txt")
	if err == nil {
		fmt.Println("\t", string(current))
	} else {
		fmt.Println("\t", "none")
	}

	fmt.Println()

	fmt.Println("Installed JAVA versions:")

	dir, err := os.ReadDir(data + "/installed")
	if err != nil {
		log.Fatal(err)
	}

	if len(dir) == 0 {
		fmt.Println("\t", "none")
	}

	for _, f := range dir {
		fmt.Println("\t", f.Name())
	}

	fmt.Println()
	fmt.Println("Available JAVA versions:")
	version := available()
	for _, v := range version {
		fmt.Println("\t", v, "\tjdk/jre")
	}
}

func use(version string) {
	fmt.Println("Changing to JDK version: " + version)

	dir, err := os.ReadDir(data + "/installed")
	if err != nil {
		log.Fatal(err)
	}

	found := false

	for _, f := range dir {
		fmt.Println(f.Name())
		if f.Name() == "jdk-"+version {
			found = true
		}
	}

	if !found {
		fmt.Println("JDK version " + version + " not found")
		fmt.Println("Use 'jem install " + version + "' to install it")
		return
	}
	
	fmt.Println("Found JDK version " + version)

	err = os.RemoveAll(data + "/current")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(data+"/current", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	link(data+"/installed/jdk-"+version, data+"/current")

	ver, err := os.Create(data + "/current_version.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = ver.WriteString("jdk-" + version)
	if err != nil {
		log.Fatal(err)
	}

	err = ver.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func install(version string) {
	fmt.Println("Installing JDK version: " + version)

	fmt.Println("Downloading...")

	temp, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	fileName := temp + "jdk-" + version + ".zip"

	download(version, osName, arch, fileName)

	fmt.Println("Extracting...")

	unzip(fileName, data+"/installed/jdk-"+version)

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(temp)

	fmt.Println("Done!")
}

func uninstall(version string) {
	fmt.Println("Uninstalling JDK version: " + version)

	fmt.Println("Removing...")
	err := os.RemoveAll(data + "/installed/jdk-" + version)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func help() {
	fmt.Println("Usage: jem [command]")
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
			if len(os.Args) > 2 {
				use(os.Args[2])
			} else {
				help()
			}
		case "install":
			if len(os.Args) > 2 {
				install(os.Args[2])
			} else {
				help()
			}
		case "uninstall":
			if len(os.Args) > 2 {
				uninstall(os.Args[2])
			} else {
				help()
			}
		default:
			help()
		}
	} else {
		help()
	}
}

func main() {
	// START TESTING ONLY
	//os.Args = append(os.Args, "list")
	//os.Args = append(os.Args, "8")
	// END TESTING ONLY

	fmt.Println("Jem - Java environment manager")
	fmt.Println("Version: " + version)
	fmt.Println()

	_ = os.Mkdir(data, os.ModePerm)
	_ = os.Mkdir(data+"/current", os.ModePerm)
	_ = os.Mkdir(data+"/installed", os.ModePerm)

	selectCommand()
}
