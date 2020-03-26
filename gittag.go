package gitcomm

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	bb "github.com/karantin2020/promptui"
)

// copypast from https://github.com/calmh/git-autotag

// levels describes git tag struct
var levels = map[string]int{
	"major": 0,
	"minor": 1,
	"patch": 2,
	"x":     0,
	"y":     1,
	"z":     2,
}

// AutoTag creates an annonated tag for the next logical version
func AutoTag(level string) {
	sign := getGitConfigBool("autotag.sign")

	curVer := closestVersion()
	if curVer == "" {
		curVer = "v0.0.0"
	}
	newVer := bumpVersion(curVer, levels[level])
	p := bb.Prompt{
		BasicPrompt: bb.BasicPrompt{
			Label:   "Need to correct tag?",
			Default: newVer,
			// Formatter: linterSubject,
			Validate: validateTag,
		},
	}
	newVer, err := p.Run()
	tagCmd := []string{"tag", "-a", "-m", "version " + newVer}
	if sign {
		tagCmd = append(tagCmd, "-s")
	}

	// fmt.Println(newVer)
	checkInterrupt(err)
	tagCmd = append(tagCmd, newVer)
	gitColorCmd(tagCmd...)
}

func getGitConfig(args ...string) string {
	args = append([]string{"config", "--get"}, args...)
	cmd := exec.Command("git", args...)
	bs, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(bs))
}

func getGitConfigBool(args ...string) bool {
	args = append([]string{"--bool"}, args...)
	return getGitConfig(args...) == "true"
}

func closestVersion() string {
	cmd := exec.Command("git", "tag")
	bs, err := cmd.Output()
	tl := bytes.Split([]byte(bs), []byte("\n"))
	if len(tl) > 1 {
		bs = tl[len(tl)-2]
	}
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(bs))
}

func bumpVersion(ver string, part int) string {
	prefix, parts := versionParts(ver)
	parts[part]++
	for i := part + 1; i < len(parts); i++ {
		parts[i] = 0
	}
	return versionString(prefix, parts)
}

func versionString(prefix string, parts []int) string {
	ver := fmt.Sprintf("%s%d", prefix, parts[0])
	for _, part := range parts[1:] {
		ver = fmt.Sprintf("%s.%d", ver, part)
	}
	return ver
}

// versionParts matches a px.y.z version, for non-digit values of p and digits
// x, y, and z.
func versionParts(s string) (prefix string, parts []int) {
	exp := regexp.MustCompile(`^([^\d]*)(\d+)\.(\d+)\.(\d+)$`)
	match := exp.FindStringSubmatch(s)
	if len(match) > 1 {
		prefix = match[1]
		parts = make([]int, len(match)-2)
		for i := range parts {
			parts[i], _ = strconv.Atoi(match[i+2])
		}
	}
	return
}

func validateTag(version string) error {
	if version == "" {
		return bb.NewValidationError("Version must not be empty string")
	}
	exp := regexp.MustCompile(`^([^\d]*)(\d+)\.(\d+)\.(\d+)$`)
	match := exp.FindStringSubmatch(version)
	if len(match) != 5 {
		return bb.NewValidationError("Version tag must be of type 'v0.1.2'")
	}
	return nil
}
