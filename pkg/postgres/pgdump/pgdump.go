package pgdump

import (
	"bytes"
	"os/exec"

	v "github.com/vilamslep/vspg/pkg/os"
)

var (
	PGDumpExe string
)

func Dump(db string, dst string, excludedTables []string) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {

	cmd := exec.Command(PGDumpExe,
		"--format", "directory", "--no-password",
		"--jobs", "4", "--blobs",
		"--encoding", "UTF8",
		"--verbose", "--file", dst,
		"--dbname", db)

	excludingArgs(cmd, excludedTables)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err = v.ExecCommand(cmd)

	return stdout, stderr, err
}

func excludingArgs(cmd *exec.Cmd, excludedTable []string) {
	for _, i := range excludedTable {
		cmd.Args = append(cmd.Args, "--exclude-table-data")
		cmd.Args = append(cmd.Args, i)
	}
}
