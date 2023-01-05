package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var data = ""

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

func use(version string, jtype string) {
	jtype = strings.TrimSpace(jtype)
	jtype = strings.ToUpper(jtype)

	fmt.Println("Changing to " + jtype + " version: " + version)

	dir, err := os.ReadDir(data + "/installed")
	if err != nil {
		log.Fatal(err)
	}

	found := false

	for _, f := range dir {
		fmt.Println(f.Name())
		if f.Name() == strings.ToLower(jtype)+"-"+version {
			found = true
		}
	}

	if !found {
		fmt.Println(jtype + " version " + version + " not found")
		fmt.Println("Use 'jem install " + version + " " + strings.ToLower(jtype) + "' to install it")
		return
	}

	fmt.Println("Found " + jtype + " version " + version)

	err = os.RemoveAll(data + "/current")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(data+"/current", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	link(data+"/installed/"+strings.ToLower(jtype)+"-"+version, data+"/current")

	ver, err := os.Create(data + "/current_version.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = ver.WriteString(strings.ToLower(jtype) + "-" + version)
	if err != nil {
		log.Fatal(err)
	}

	err = ver.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func install(version string, jtype string) {
	jtype = strings.TrimSpace(jtype)
	jtype = strings.ToUpper(jtype)

	fmt.Println("Installing " + jtype + " version: " + version)

	fmt.Println("Downloading...")

	temp, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	fileName := temp + strings.ToLower(jtype) + "-" + version + ".zip"

	download(version, jtype, osName, arch, fileName)

	fmt.Println("Extracting...")

	unzip(fileName, data+"/installed/"+strings.ToLower(jtype)+"-"+version)

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(temp)

	fmt.Println("Done!")
}

func uninstall(version string, jtype string) {
	jtype = strings.TrimSpace(jtype)
	jtype = strings.ToUpper(jtype)

	fmt.Println("Uninstalling " + jtype + " version: " + version)

	fmt.Println("Removing...")
	err := os.RemoveAll(data + "/installed/" + strings.ToLower(jtype) + "-" + version)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func help() {
	fmt.Println("Usage: jem [command]")
	fmt.Println("Commands:")
	fmt.Println("  list")
	fmt.Println("  use [version] [jre|jdk] (default jdk)")
	fmt.Println("  install [version] [jre|jdk] (default jdk)")
	fmt.Println("  uninstall [version] [jre|jdk] (default jdk)")
}
func selectCommand() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			list()
		case "use":
			if len(os.Args) > 3 {
				use(os.Args[2], os.Args[3])
			} else if len(os.Args) > 2 {
				use(os.Args[2], "jdk")
			} else {
				help()
			}
		case "install":
			if len(os.Args) > 3 {
				install(os.Args[2], os.Args[3])
			} else if len(os.Args) > 2 {
				install(os.Args[2], "jdk")
			} else {
				help()
			}
		case "uninstall":
			if len(os.Args) > 3 {
				uninstall(os.Args[2], os.Args[3])
			} else if len(os.Args) > 2 {
				uninstall(os.Args[2], "jdk")
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

	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	data = filepath.Dir(exec)

	//fmt.Println("Data directory: " + data)

	_ = os.Mkdir(data, os.ModePerm)
	_ = os.Mkdir(data+"/current", os.ModePerm)
	_ = os.Mkdir(data+"/installed", os.ModePerm)

	selectCommand()
}
