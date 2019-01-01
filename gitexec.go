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
	cmd := exec.Command("git", "-c", "color.ui=always", "commit", "-a", "-m", msg)
	log.Printf("git commit -a -m \"%s\"\n", msg)
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	fmt.Println(string(out))
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
	log.Printf("git status --porcelain\n")
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	return len(out) != 0
}
