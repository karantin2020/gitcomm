package gitcomm

import (
	"fmt"
	"log"
	"os/exec"
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
	aflag := "-a"
	if addAll {
		aflag = ""
	}
	cmd := exec.Command("git", "-c", "color.ui=always", "commit", aflag, "-m", msg)
	log.Printf("git commit %s -m \"%s\"\n", aflag, msg)
	out, err := cmd.CombinedOutput()
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

// CheckForUncommited function checkes if there are changes that need commit
func CheckForUncommited() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	// log.Printf("git status --porcelain\n")
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	return len(out) != 0
}
