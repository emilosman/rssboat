package rss

import "errors"

var (
	ErrFeedHasNoUrl       = errors.New("Feed has no URL")
	ErrNoFeedsInList      = errors.New("No feeds in list")
	ErrNoCategoryGiven    = errors.New("No category given")
	ErrCacheEmpty         = errors.New("Cache empty")
	ErrConfigDoesNotExist = "open urls.yaml: file does not exist"
	MsgFeedNotLoaded      = "Feed not loaded yet. Press shift+r"
)
