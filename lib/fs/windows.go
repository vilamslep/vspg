package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/vilamslep/psql.maintenance/lib/config"
)

func IsEnoughSpace(src string, dst string) (bool, error) {

	free, err := freeSpace(dst)

	if err != nil {
		return false, err
	}

	used, err := GetSize(src)

	if err != nil {
		return false, err
	}
	return free > used, err
}

func freeSpace(path string) (int64, error) {
	kernelDLL := syscall.MustLoadDLL("kernel32.dll")
	GetDiskFreeSpaceExW := kernelDLL.MustFindProc("GetDiskFreeSpaceExW")

	var free int64

	r1, _, err := GetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
	)

	if r1 == 0 {
		return 0, fmt.Errorf("%s. Returned code %d", err.Error(), r1)
	}

	return free, nil
}

func GetRootDir(path string, name string, kind int) (string, error) {
	now := time.Now()

	var ft string
	switch kind {
	case config.DAILY:
		ft = now.Format("02-01-2006")
	case config.WEEKLY:
		_, week := now.ISOWeek()
		ft = strconv.Itoa(week)
	case config.MONTHLY:
		ft = now.Format("01")
	default:
		return "", fmt.Errorf("undefined kind. kind %d", kind)
	}

	backPath := fmt.Sprintf("%s\\%s\\%s", path, name, ft)

	return backPath, createIfNotExists(backPath)
}

func CreateDirectories(root string, name string, children []string) (location map[string]string, err error) {
	path := fmt.Sprintf("%s\\%s", root, name)
	if err = createIfNotExists(path); err != nil {
		return
	}

	location["main"] = path

	for _, ch := range children {
		chPath := fmt.Sprintf("%s\\%s", path, ch)
		if err := createIfNotExists(chPath); err != nil {
			return nil, err
		}
		location[ch] = chPath
	}
	return
}

func ClearOldBackup(path string, count int) (err error) {

	isDir, err := isDir(path)
	if err != nil {
		return
	}

	if !isDir {
		return fmt.Errorf("%s isn't directory", path)
	}

	if ls, err := ioutil.ReadDir(path); err == nil {
		sort.Slice(ls, func(i, j int) bool {
			return ls[i].ModTime().Before(ls[j].ModTime())
		})
		if len(ls) < count {
			return err
		}
		toRemove := ls[0:(len(ls) - count)]
		return removeFile(path, toRemove)

	} else {
		return err
	}
}

func removeFile(path string, files []os.FileInfo) error {
	for _, fi := range files {
		fpath := fmt.Sprintf("%s\\%s", path, fi.Name())
		if err := os.RemoveAll(fpath); err != nil {
			return err
		}
	}
	return nil
}
