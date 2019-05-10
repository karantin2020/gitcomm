package gitcomm

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
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
	tagCmd := []string{"tag", "-a", "-m", "version " + newVer}
	if sign {
		tagCmd = append(tagCmd, "-s")
	}
	tagCmd = append(tagCmd, newVer)

	fmt.Println(newVer)
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
	bs = tl[len(tl)-2]
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
