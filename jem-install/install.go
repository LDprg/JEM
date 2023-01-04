package main

import (
	"fmt"
	"log"
	"os"
)

func copyFile(src string, dest string) {
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

func main() {
	fmt.Println("Installing JEM...")

	// PATH Not Working
	_ = os.MkdirAll("~/.jem", os.ModePerm)

	// to prevent path bugs
	_ = os.MkdirAll("~/.jem/current/bin", os.ModePerm)

	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		copyFile(f.Name(), "~/.jem/"+f.Name())
	}

	fmt.Println("Done!")
}
