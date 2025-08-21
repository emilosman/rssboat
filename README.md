# rssboat
- RSS reader inspired by newsboat

## MVP
- [ ] Update JSON items via GUID
- [ ] UpdateAll() should preserve read state

## Todo
- [ ] Messages moved to vars
- [ ] Unit test tui funcs (BuildFeedList...)
- [ ] Timeout network request 10s
- [ ] Sort feeds
- [ ] Sort feed items by date
- [ ] Unread items should have different color
- [ ] Feeds should load async on UpdateAll()
- [ ] Should send update command after go func finishes
- [ ] Feed items should have a view that displays content
- [ ] feeds.yaml should not be relative path
- [ ] "o" should not trigger feed open when filtering
- [ ] Unset filter state on back navigation
- [ ] AdditionalShortHelpKeys() extend
- [ ] AdditionalFullHelpKeys() extend

## Refactor
- [ ] DRY cleanup
- [ ] Refactor use of pointers
- [ ] Alphabetize key commands in Update func
- [ ] Refactor list updates into helper func
- [ ] Refactor view selection m.selectedFeed
- [ ] Refactor feed/item open error messages and funcs

## Future
- [ ] E2E tests [teatest](https://github.com/caarlos0/teatest-example/blob/main/main_test.go)
- [ ] Opening feed items should work on all operating systems with default browser
- [ ] Public repo...

## Maybe
- [ ] Confirmation Y/N on major commands
- [ ] Tabs for feed categories
- [ ] "h" and "l" should open and close feeds
- [ ] Sort options (1. newest unread up top, 2. popular...)

## Done
- [x] "C" should mark entire feed list read
- [x] Fix item pointers and state persistance
- [x] Toggle read state on item open
- [x] Save/restore on app open
- [x] "q" should not quit when in items view, but should go back to feeds view
- [x] State should be stored in JSON
- [x] "a" should toggle item read state manually
- [x] "q" quit disabled
- [x] "A" should mark feed as read
- [x] "r" should refresh a single feed
- [x] "Latest" should be displayed in feed item description when present
- [x] Make field test names descriptive
- [x] Table test field names
- [x] Test unloadedFeed
