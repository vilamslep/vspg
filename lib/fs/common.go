package fs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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

func Remove(paths ...string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		} 
	}
	return nil	
}

func TempDir() (string, error) {
	err := CreateIfNotExists(WIN_OS_PROGDATA)
	if err != nil {
		return "", err
	}
	return WIN_OS_PROGDATA, nil
}

func CreateIfNotExists(path string) error {
	return os.MkdirAll(path, 0777)
}

func LoadEnvfile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		path := strings.Split(sc.Text(), "=")
		if len(path) < 2 {
			continue
		}
		os.Setenv(path[0], path[1])
	}
	return sc.Err()
}
