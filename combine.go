package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"

	"github.com/spf13/cobra"
)

var (
	build string
	keep  bool
	note  string
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

		if len(build) == 0 {
			fmt.Println("Build is required see --help for usage")
			return
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
		header := build + " (" + strings.Fields(time.Now().String())[0] + ")"
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "## %s\n", header)
		if note != "" {
			fmt.Fprintf(b, "**%s**\n", note)
		}
		for _, e := range changelogEntries {
			result := "* "

			if e.Platform != "" {
				result += "[" + e.Platform + "] "
			}
			result += e.Message + " (" + e.Author + ")"

			if e.Merge != "" {
				result += " !" + e.Merge + "\n"
			}
			fmt.Fprint(b, result)
		}

		fmt.Fprintf(b, "\n")

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

		ioutil.WriteFile(rootPath+build+"-CHANGELOG.md", b.Bytes(), 0644)

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
}
