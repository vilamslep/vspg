package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	WIN_OS_PROGDATA = "C:\\Temp\\postgres.backup"
)

func GetSize(path string) (int64, error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		info, err := file.Stat()
		if err != nil {
			return 0, err
		}
		if !info.IsDir() {
			return info.Size(), err
		}
	} else {
		return 0, err
	}
	return getDirectorySize(path)
}

func isDir(path string) (rs bool, err error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		info, err := file.Stat()
		if err != nil {
			return rs, err
		}
		return info.IsDir(), err
	} else {
		return rs, err
	}
}

func getDirectorySize(path string) (totalSize int64, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func TempDir() (string, error) {
	err := createIfNotExists(WIN_OS_PROGDATA)
	if err != nil {
		return "", err
	}
	return WIN_OS_PROGDATA, nil
}

func createIfNotExists(path string) error {
	return os.MkdirAll(path, 0777)
}

func Copy(src string, dst string) error {

	fdr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fdr.Close()

	stat, err := fdr.Stat()

	if err != nil {
		return err
	}
	isDir := stat.IsDir()

	if !isDir {
		if err := CopyFile(fdr, dst); err != nil {
			return err
		}
	} else {
		if err := CopyDirectory(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func CopyFile(fdr *os.File, dst string) error {
	fdw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdw.Close()

	if _, err := io.Copy(fdw, fdr); err != nil {
		return err
	}
	return nil
}

func CopyDirectory(src string, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := createIfNotExists(dstPath); err != nil {
				return nil
			}
		}
		
		if err := Copy(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}
