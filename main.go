package main

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var (
	rootPath             string
	changelogEntriesPath string
	changelogFilePath    string
)

const version = "0.1"

// ChangelogEntry as a databucket for log entry data
type ChangelogEntry struct {
	Message  string `json:"message"`
	Author   string `json:"author"`
	Platform string `json:"platform"`
	Merge    string `json:"merge"`
}

func main() {
	path, _ := getAppPath()
	rootPath = path
	changelogEntriesPath = path + "changelogs/unreleased/"
	changelogFilePath = path + "CHANGELOG.md"

	Execute()
}

func getAppPath() (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "--show-toplevel")

	b := &bytes.Buffer{}
	gitCmd.Stdout = b

	err := gitCmd.Run()
	if err != nil {
		return "", errors.Wrap(err, "cannot find app root dir git rev-parse failed")
	}

	output := b.String()

	if len(output) == 0 {
		return "", errors.New("cannot find app root dir git rev-parse had no output")
	}

	return strings.TrimSpace(output) + "/", nil
}
