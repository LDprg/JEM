package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("Uninstalling JEM...")

	data, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	jemPath := data + "/.jem"

	jemPath = filepath.Clean(jemPath)

	fmt.Println("Removing...")

	err = os.RemoveAll(jemPath)
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell", "[Environment]::GetEnvironmentVariable(\"Path\",\"User\")").Output()
		if err != nil {
			log.Fatal(err)
		}

		path := string(out)

		newPath := ""

		for _, p := range strings.Split(path, ";") {
			p = strings.TrimSpace(p)
			if !strings.Contains(p, ".jem") && p != "" {
				newPath += p + ";"
			}
		}

		err = exec.Command("powershell", "[Environment]::SetEnvironmentVariable(\"Path\",\""+newPath+"\",\"User\")").Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := exec.Command("sed -in \"/\\.jem/d\" $HOME/.bashrc").Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done!")
}
