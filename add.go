package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	author   string
	merge    string
	message  string
	platform string
	entry    ChangelogEntry
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new changelog entry file",
	Long: `Adds a changelog entry file requiring description and author. Example usage:
	
	add -a="John Smith" -d="full description text here" -m="123"`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(author) == 0 {
			fmt.Println("Author is required see --help for usage")
			return
		} else if len(message) == 0 {
			fmt.Println("Description is required see --help for usage")
			return
		}

		entry.Author = author
		entry.Platform = strings.ToLower(platform)
		entry.Message = message
		entry.Merge = merge

		entryJSON, err := json.Marshal(entry)
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := os.Stat(changelogEntriesPath); os.IsNotExist(err) {
			pathErr := os.MkdirAll(changelogEntriesPath, 0777)
			if pathErr != nil {
				fmt.Println(err)
				return
			}
		}
		filename := strings.Replace(strings.Fields(time.Now().String())[0], "-", "", -1)
		filename += "-" + sanitizeDescription(entry.Message) + ".json"
		err = ioutil.WriteFile(changelogEntriesPath+filename, entryJSON, 0644)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&author, "author", "a", "", "Required: name of the user submitting the message")
	addCmd.Flags().StringVarP(&message, "description", "d", "", "Required: description to summarize changes")
	addCmd.Flags().StringVarP(&merge, "merge", "m", "", "Merge request id")
	addCmd.Flags().StringVarP(&platform, "platform", "p", "", "Platform affected by changes")
}

func sanitizeDescription(msg string) string {
	if msg == "" {
		return msg
	}
	result := msg
	invalidChars := []string{"*", "`", "_", "~", "]", "[", "!", "#", "(", ")"}

	for _, c := range invalidChars {
		result = strings.Replace(result, c, "", -1)
	}

	result = strings.Replace(result, " ", "-", -1)

	return result
}
