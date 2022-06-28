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
	if _, err := IsEnoughSpace("d:\\gtp", "c:\\"); err != nil {
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
