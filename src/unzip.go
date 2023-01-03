package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func unzip(file string, dest string) {
	if strings.Contains(file, ".zip") {

		archive, err := zip.OpenReader(file)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range archive.File {
			filePath := filepath.Join(dest, f.Name)
			fmt.Println("unzipping file ", filePath)

			if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return
			}
			if f.FileInfo().IsDir() {
				fmt.Println("creating directory...")
				_ = os.MkdirAll(filePath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				panic(err)
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(err)
			}

			fileInArchive, err := f.Open()
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				panic(err)
			}

			_ = dstFile.Close()
			_ = fileInArchive.Close()
		}

	} else {
		log.Fatal("Invalid file type")
	}

}
