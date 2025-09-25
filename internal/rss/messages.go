package rss

import "errors"

var (
	ErrFeedHasNoUrl    = errors.New("Feed has no URL")
	ErrNoFeedsInList   = errors.New("No feeds in list")
	ErrNoCategoryGiven = errors.New("No category given")
	ErrChacheEmpty     = errors.New("Cache empty")
	MsgFeedNotLoaded   = "Feed not loaded yet. Press shift+r"
)
