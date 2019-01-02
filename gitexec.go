package gitcomm

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// GitExec function performs git workflow
func GitExec(addAll, show bool, msg string) {
	if addAll {
		cmd := exec.Command("git", "-c", "color.ui=always", "add", "-A")
		log.Printf("git add -A\n")
		out, err := cmd.CombinedOutput()
		CheckIfError(err)
		fmt.Println(string(out))
	}
	cmd := exec.Command("git", "add", "-u")
	log.Printf("git add -u")
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	fmt.Printf("\n%s\n", string(out))

	cmd = exec.Command("git", "-c", "color.ui=always", "commit", "-m", msg)
	log.Printf("git commit -m\n")
	out, err = cmd.CombinedOutput()
	CheckIfError(err)
	fmt.Printf("\n%s\n", string(out))
	if show {
		cmd := exec.Command("git", "-c", "color.ui=always", "show", "-s")
		log.Printf("git show -s\n")
		out, err := cmd.CombinedOutput()
		CheckIfError(err)
		fmt.Println(string(out))
	}
}

// CheckForUncommited function checks if there are changes that need commit
func CheckForUncommited() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	return len(out) != 0
}

// CheckIsGitDir function checks is dir inside git worktree
func CheckIsGitDir() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	isGitDir := strings.TrimSpace(string(out))
	if err == nil && isGitDir == "true" {
		return true
	}
	return false
}
