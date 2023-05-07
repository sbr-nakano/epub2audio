package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func unzipEpub(epubPath, tmpPath string) error {
	zipReader, err := zip.OpenReader(epubPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	xhtmlPath := filepath.Join(tmpPath, "xhtml")
	err = os.MkdirAll(xhtmlPath, 0755)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		ext := filepath.Ext(file.Name)
		if ext != ".xhtml" && ext != ".html" {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		targetPath := filepath.Join(xhtmlPath, filepath.Base(file.Name))
		targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
