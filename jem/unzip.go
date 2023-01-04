package main

import (
	"archive/zip"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func link(source string, dest string, isbar ...bool) {
	source = filepath.Clean(source) + string(os.PathSeparator)
	dest = filepath.Clean(dest) + string(os.PathSeparator)

	src, err := os.ReadDir(source)
	if err != nil {
		log.Fatal(err)
	}

	bar := pb.New(len(src))

	if len(isbar) == 0 {
		bar.Start()
	}

	for _, f := range src {
		if len(isbar) == 0 {
			bar.Increment()
		}

		if f.IsDir() {
			err = os.Mkdir(dest+f.Name(), os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			link(source+f.Name(), dest+f.Name(), false)
		} else {
			err = os.Link(source+f.Name(), dest+f.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if len(isbar) == 0 {
		bar.Finish()
	}
}

func unzip(file string, dest string) {
	if strings.Contains(file, ".zip") {

		err := os.RemoveAll(dest)
		if err != nil {
			log.Fatal(err)
		}

		archive, err := zip.OpenReader(file)
		if err != nil {
			log.Fatal(err)
		}

		count := len(archive.File)
		bar := pb.StartNew(count)

		for _, f := range archive.File {
			bar.Increment()

			_, fileName, _ := strings.Cut(f.Name, "/")

			if strings.Trim(fileName, " ") == "" {
				continue
			}

			filePath := filepath.Join(dest, fileName)

			if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return
			}
			if f.FileInfo().IsDir() {
				_ = os.MkdirAll(filePath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				log.Fatal(err)
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatal(err)
			}

			fileInArchive, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				log.Fatal(err)
			}

			_ = dstFile.Close()
			_ = fileInArchive.Close()
		}

		bar.Finish()
	} else {
		log.Fatal("Invalid file type")
	}
}
