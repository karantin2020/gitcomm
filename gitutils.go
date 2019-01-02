package gitcomm

import (
	"os"
	"os/exec"
)

var (
	colorOutputFlag = []string{"-c", "color.ui=always"}
)

// UndoLastCommit leaves your working tree (the state of your
// files on disk) unchanged but undoes the commit
// and leaves the changes you committed unstaged
func UndoLastCommit() {
	undoCmd := []string{"reset", "HEAD~"}
	args := append(colorOutputFlag, undoCmd...)
	git(args...)
}

func git(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	ExitIfError(err)
}
