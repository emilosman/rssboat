package rss

import "errors"

var (
	ErrFeedHasNoUrl       = errors.New("Feed has no URL")
	ErrNoFeedsInList      = errors.New("No feeds in list")
	ErrNoCategoryGiven    = errors.New("No category given")
	ErrConfigDoesNotExist = "open urls.yaml: file does not exist"
	MsgFeedNotLoaded      = "Feed not loaded yet. Press shift+r"
	ExampleConfigFile     = `# This file is written in YAML format.
# Each feed must be organized under a category.
# Feeds that are not assigned to a category will NOT appear in the app.
# Formatting Rules:
# - Use two spaces for indentation (no tabs allowed).
# - Follow proper YAML structure (see example below).
# Additional Information:
# - Refer to the README.md file for details on where this configuration file is stored
# on your operating system.
# - After saving and closing this file, the app will automatically load and display your feeds.
# Enjoy!
#
# Example (uncomment lines below to use):
#golang:
#  - https://go.dev/blog/feed.atom
#  - https://golang.cafe/rss
`
)
