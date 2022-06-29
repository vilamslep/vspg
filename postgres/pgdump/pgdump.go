package pgdump

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
)
var(
	PGDump string
)

func Dump(db string, dst string, output io.Writer, excludedTables []string) error {

	cmd := exec.Command(PGDump, "--format", "directory", "--no-password","--jobs", "4",
		"--blobs", "--encoding", "UTF8", "--verbose","--file", dst, "--dbname", db)
	excludingArgs(cmd, excludedTables)

	if err := cmd.Start(); err != nil {
        return err
    }
	cmd.Stdout = output
    if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                return fmt.Errorf("Exit Status: %d", status.ExitStatus())
            }
        } else {
            return fmt.Errorf("cmd.Wait: %v", err)
        }
    }
	return nil
}

func excludingArgs(cmd *exec.Cmd, excludedTable []string) {
	for _, i := range excludedTable {
		cmd.Args = append(cmd.Args, "--exclude-table-data")
		cmd.Args = append(cmd.Args, i)
	}
}
	
