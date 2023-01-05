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

func copyFile(src string, dest string) {
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)
	fmt.Println("Copying " + src + " to " + dest)

	content, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(dest, content, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func addPath(newpath string) {
	newpath = filepath.Clean(newpath)

	if runtime.GOOS == "windows" {
		err := exec.Command("powershell", "[Environment]::SetEnvironmentVariable(\"Path\",[Environment]::GetEnvironmentVariable(\"Path\",\"User\") + \";\"+(Resolve-Path ", newpath, "), \"User\")\n").Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		exec.Command("echo export PATH=", newpath, ":\\$PATH >> $HOME/.bashrc")
	}
}

func main() {
	fmt.Println("Installing JEM...")

	data, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// PATH Not Working
	_ = os.MkdirAll(data+"/.jem", os.ModePerm)

	// to prevent path bugs
	_ = os.MkdirAll(data+"/.jem/current/bin", os.ModePerm)

	// ADD JEM TO PATH
	jemPath := data + "/.jem"
	javaPath := data + "/.jem/current/bin"

	jemFound := false
	javaFound := false

	jemPath = filepath.Clean(jemPath)
	javaPath = filepath.Clean(javaPath)

	path := os.Getenv("PATH")

	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell", "[Environment]::GetEnvironmentVariable(\"Path\",\"User\")").Output()
		if err != nil {
			return
		}

		path = string(out)
	}

	for _, p := range strings.Split(string(path), ";") {
		p = filepath.Clean(p)

		if strings.Contains(p, jemPath) {
			jemFound = true
		}
		if strings.Contains(p, javaPath) {
			javaFound = true
		}
	}

	if jemFound {
		fmt.Println("JEM already in PATH")
	} else {
		fmt.Println("JEM add to PATH")
		addPath(jemPath)
	}

	if javaFound {
		fmt.Println("JEM JAVA already in PATH")
	} else {
		fmt.Println("JEM JAVA add to PATH")
		addPath(javaPath)
	}

	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		copyFile(f.Name(), data+"/.jem/"+f.Name())
	}

	fmt.Println("Done!")
}
