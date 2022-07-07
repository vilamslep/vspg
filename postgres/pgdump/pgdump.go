package pgdump

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

var (
	PGDumpExe string
)

func Dump(db string, dst string, excludedTables []string) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {

	cmd := exec.Command(PGDumpExe, "--format", "directory", "--no-password", "--jobs", "4",
		"--blobs", "--encoding", "UTF8", "--verbose", "--file", dst, "--dbname", db)
	excludingArgs(cmd, excludedTables)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		return stdout, stderr, err
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return stdout, stderr, fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			return stdout, stderr, fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	return stdout, stderr, nil
}

func excludingArgs(cmd *exec.Cmd, excludedTable []string) {
	for _, i := range excludedTable {
		cmd.Args = append(cmd.Args, "--exclude-table-data")
		cmd.Args = append(cmd.Args, i)
	}
}
