package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/spf13/cobra"
)

var (
	build   string
	keep    bool
	note    string
	flavor  string
	archive bool
)

// combineCmd represents the combine command
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "Combines existing changelog entries into CHANGELOG.md",
	Long: `combines all the .json files in the /changelog/unreleased/ folder
	and appends to a master CHANGELOG.md 
	as well as creating a [build]-CHANGELOG.md. Example usage:
	
	combine -b "v0.1" -n "released to app store"`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(changelogEntriesPath); os.IsNotExist(err) {
			fmt.Println(err)
			return
		}

		if _, err := os.Stat(changelogArchivePath); os.IsNotExist(err) {
			pathErr := os.MkdirAll(changelogArchivePath, 0777)
			if pathErr != nil {
				fmt.Println(err)
				return
			}
		}

		if len(build) == 0 {
			fmt.Println("Build is required see --help for usage")
			return
		}
		if archive {
			if _, err := os.Stat(changelogArchivePath + sanitizeDescription(build) + "/"); os.IsNotExist(err) {
				pathErr := os.MkdirAll(changelogArchivePath+sanitizeDescription(build)+"/", 0777)
				if pathErr != nil {
					fmt.Println(err)
					return
				}
			}
		}

		//Parse files
		var changelogFiles []os.FileInfo
		var changelogEntries []ChangelogEntry

		files, err := ioutil.ReadDir(changelogEntriesPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" {
				changelogFiles = append(changelogFiles, f)

				contents, err := ioutil.ReadFile(changelogEntriesPath + f.Name())
				if err != nil {
					fmt.Println(err)
					return
				}

				//archive old changelogFiles
				if archive {
					err = ioutil.WriteFile(changelogArchivePath+sanitizeDescription(build)+"/"+f.Name(), contents, 0644)
					if err != nil {
						fmt.Println(err)
					}
				}

				var entry ChangelogEntry
				err = json.Unmarshal(contents, &entry)
				if err != nil {
					fmt.Println(err)
					return
				}

				changelogEntries = append(changelogEntries, entry)
			}
		}

		//Build output
		b := &bytes.Buffer{}
		if len(flavor) == 0 {
			unflavoredContent(b, changelogEntries)
		} else {
			buildFlavoredContent(strings.ToLower(flavor), b, changelogEntries)
		}

		//Write files
		os.OpenFile(changelogFilePath, os.O_RDONLY|os.O_CREATE, 0666)

		fileContent, err := ioutil.ReadFile(changelogFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		ioutil.WriteFile(changelogFilePath, b.Bytes(), 0644)

		f, err := os.OpenFile(changelogFilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer f.Close()

		if _, err = f.Write(fileContent); err != nil {
			fmt.Println(err)
			return
		}

		ioutil.WriteFile(changelogArchivePath+build+"-CHANGELOG.md", b.Bytes(), 0644)

		//remove old changelogFiles
		if !keep {
			for _, f := range changelogFiles {
				os.Remove(changelogEntriesPath + f.Name())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(combineCmd)
	combineCmd.Flags().StringVarP(&build, "build", "b", "", "Required: Build version used as header in .md file")
	combineCmd.Flags().StringVarP(&note, "note", "n", "", "Note about this build")
	combineCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Maintains the .json files in the unreleased directory")
	combineCmd.Flags().StringVarP(&flavor, "flavor", "f", "", "Sets markdown flavor, can be 'github' or 'gitlab'")
	combineCmd.Flags().BoolVarP(&archive, "archive", "a", false, "Archives the .json files to changelogs/released/[build]/")
}
