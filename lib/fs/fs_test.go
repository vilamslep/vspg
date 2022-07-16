package fs

import (
	"testing"
)

func TestGetSize(t *testing.T) {
	if size, err := GetSize("d:\\gtp"); err == nil {
		if size != 24202416 {
			t.Fatalf("wrong count. result %d. Expected %d", size, 24202416)
		}
	} else {
		t.Fatal(err)
	}
}

func TestIsEnoughSpace(t *testing.T) {
	if _, err := IsEnoughSpace("d:\\gtp", "c:\\", 0); err != nil {
		t.Fatal(err)
	}
}

func TestClearOldBackup(t *testing.T) {
	if err := ClearOldBackup("d:\\gtp", 1); err != nil {
		t.Fatal(err)
	}
}

func TestCompress(t *testing.T) {
	if err := Compress("d:\\gtp", "d:\\gtp.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {
	if err := Copy("C:\\Pictures\\Art\\2.jpg", "C:\\backup\\2.jpg"); err != nil {
		t.Fatal(err)
	}
}

func TestCopyDirectory(t *testing.T) {
	if err := Copy("C:\\backup\\daily\\kfk", "C:\\backup\\kfk"); err != nil {
		t.Fatal(err)
	}
}
