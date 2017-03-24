package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func buildFlavoredContent(flavor string, b *bytes.Buffer, changelogEntries []ChangelogEntry) {
	mergeIdentifier := "#"
	if flavor == "gitlab" {
		mergeIdentifier = "!"
	}
	header := build + " (" + strings.Fields(time.Now().String())[0] + ")"
	fmt.Fprintf(b, "## %s\n", header)
	if note != "" {
		fmt.Fprintf(b, "**%s**\n", note)
	}
	for _, e := range changelogEntries {
		result := "* "

		if e.Platform != "" {
			result += "[" + e.Platform + "] "
		}
		result += e.Message + " @" + e.Author

		if e.Merge != "" {
			result += " " + mergeIdentifier + e.Merge + "\n"
		}
		fmt.Fprint(b, result)
	}

	fmt.Fprintf(b, "\n")
}

func unflavoredContent(b *bytes.Buffer, changelogEntries []ChangelogEntry) {
	header := build + " (" + strings.Fields(time.Now().String())[0] + ")"
	fmt.Fprintf(b, "## %s\n", header)
	if note != "" {
		fmt.Fprintf(b, "**%s**\n", note)
	}
	for _, e := range changelogEntries {
		result := "* "

		if e.Platform != "" {
			result += "[" + e.Platform + "] "
		}
		result += e.Message + " - *" + e.Author + "*"

		if e.Merge != "" {
			result += " - **Merge:** " + e.Merge + "\n"
		}
		fmt.Fprint(b, result)
	}

	fmt.Fprintf(b, "\n")
}
