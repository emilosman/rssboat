# rssboat
- RSS reader inspired by newsboat

## MVP
- [ ] Feeds should load async on UpdateAll(), update one-by-one
- [ ] Refresh status messages fix

## Todo
- [ ] Increase test coverage
- [ ] Unread counter (15/254)
- [ ] Messages moved to vars
- [ ] Unit test tui funcs (BuildFeedList...)
- [ ] Timeout network request 10s
- [ ] Unread items should have different color
- [ ] Should send update command after go func finishes
- [ ] Feed items should have a view that displays content
- [ ] feeds.yml should not be relative path
- [ ] AdditionalShortHelpKeys() extend
- [ ] AdditionalFullHelpKeys() extend
- [ ] Reset selected item index when opening feed, remember it on the feedlist

## Refactor
- [ ] DRY cleanup
- [ ] Refactor list updates into helper func
- [ ] Refactor feed/item open error messages and funcs
- [ ] Refactor use of pointers ?

## Future
- [ ] E2E tests [teatest](https://github.com/caarlos0/teatest-example/blob/main/main_test.go)
- [ ] Public repo
- [ ] GH sponsors

## Maybe
- [ ] Confirmation Y/N on major commands
- [ ] Tabs for categories
- [ ] "h" and "l" should open and close feeds
- [ ] Sort options (1. newest unread up top, 2. popular...)

## Done
- [x] Sort feed items by date
- [x] "o" should not trigger feed open when filtering
- [x] Disable key handlers when filtering
- [x] Unset filter state on back navigation
- [x] Mark read instead of toggle on open
- [x] Update JSON items via GUID
- [x] UpdateAll() should preserve read state
- [x] Opening feed items should work on all operating systems with default browser
- [x] Refactor view selection m.selectedFeed
- [x] Alphabetize key commands in Update func
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
