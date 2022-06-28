package fs

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(src string, dst string) error {
	zipfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(src)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	filepath.Walk(src, treeHandler(archive, baseDir, src))

	return err
}

func treeHandler(archive *zip.Writer, baseDir string, src string) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, src))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		if writer, err := archive.CreateHeader(header); err == nil {
			if info.IsDir() {
				return nil
			}
			return addFileToArchive(path, writer)
		} else {
			return err
		}
	}
}

func addFileToArchive(path string, writer io.Writer) error {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		_, err = io.Copy(writer, file)

		return err
	} else {
		return err
	}
}
