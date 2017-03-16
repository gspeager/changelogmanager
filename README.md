# Changelog Manager
A small command line utility to manage changelogs and automate generating CHANGELOG.md

## Overview

The compiled Changelog Manager executable and just be dropped into any development project, like a `[project root]/bin` or `[project root]/tools` directory.  It can be used to generate json files representing a change log entry and used to combine all the files into your project's CHANGELOG.md as well as generate a standalone [build]-CHANGELOG.md in the project root directory.

> note: Changelog Manger currently uses `git rev-parse --show-toplevel` to set `[project root]`

## Usage

Changelog Manager has four commands:
 * `version` - gets the current changelogmanager version
 * `help` - displays the terminal help for changelogmanager
 * `add` - generates changelog entry .json files
 * `combine` - combines the existing changelog entry files and outputs two files
> note: all commands support -h or --help thanks to [Cobra](https://github.com/spf13/cobra)
 
### add
**add** is used to output generated .json files to `[project root]/changelogs/unreleased/`.  This command has four flags that can be set when running executing it.
* Required flags:
    * `--description [string]` or `-d [string]` - the main message that will be displayed in the changelog for the corresponding code changes.
    * `--author [string]` or `-a [string]` - name of the developer responsible for the changes
* Optional flags:
    * `--merge [id]` or `-m [id]` - sets the merge or pull request id associated with the changes
    * `--platform [string]` or `-p [string]` - logs the platform the changes were applied to
    
#### Example Usage
`changelogmanager add -d "This is my change description" -a "John Smith" -m 123`

### combine
**combine** takes all the .json file in `[project root]/changelogs/unreleased/`, combines them into a chunk of markdown, then appends the resulting markdown to the top of your `[project root]/CHANGELOG.md` as well as generating a `[project root]/[build]-CHANGELOG.md` stand alone file.  By default it will delete the .json files after updating CHANGELOG.md. Combine only has one required flag:
* Required flags:
    * `--build [string]` or `-b [string]` - sets the build/version number all the changes are going to be applied to.
* Optional flags:
    * `--note [string]` or `-n [string]` - sets a note for the build
    * `--keep` or `-k` - retains the .json files after the combine is complete
    
#### Example Usage
`changelogmanager combine -b "v0.1" -note "first build hosted on gitlab" -k`
