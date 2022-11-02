package pgdump

import (
	"bufio"
	"os"
	"testing"
)

func TestDump(t *testing.T) {
	PGDumpExe = "C:\\PostgreSQL\\bin\\pg_dump.exe"
	os.Setenv("PGUSER", "postgres")
	os.Setenv("PGPASSWORD", "142543")

	f, _ := os.Create("kfk.log")
	defer f.Close()

	stdout, stderr, err := Dump("kfk", "C:\\ut\\logic", []string{"public._inforg12487", "public.config"})

	if err != nil {
		wr := bufio.NewWriter(f)
		wr.Write(stderr.Bytes())
		wr.Flush()
		
		t.Fatal(err)
	}
	wr := bufio.NewWriter(f)
	if stdout.Len() > 0 {
		wr.Write(stdout.Bytes())
	} 
	
	if stderr.Len() > 0 {
		wr.Write(stderr.Bytes())
	}
	wr.Flush()
	t.Log("here")
}
