package gitcomm

import (
	"log"
	"os"
	"os/exec"
	"strings"
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

func gitColorCmd(cmd ...string) {
	if cmd == nil || len(cmd) == 0 {
		return
	}
	args := append(colorOutputFlag, cmd...)
	log.Printf("git %s\n", strings.Join(cmd, " "))
	git(args...)
}
