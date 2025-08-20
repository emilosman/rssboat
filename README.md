# rssboat
- RSS reader inspired by newsboat

## Todo
- [ ] Save/restore on app open
- [ ] Update JSON items via GUID
- [ ] Fix item pointers and state persistance
- [ ] Unit test tui funcs
- [ ] Timeout network request 10s
- [ ] Toggle read state on item open
- [ ] "C" should mark entire feed list read
- [ ] Sort feeds
- [ ] Sort feed items by date
- [ ] UpdateAll() called multiple times should update the list
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
- [ ] Refactor list updates into helper func
- [ ] Refactor view selection m.selectedFeed
- [ ] Refactor feed/item open error messages and funcs

## Future
- [ ] Opening feed items should work on all operating systems with default browser

## Maybe
- [ ] Tabs for feed categories
- [ ] "q" should not quit when in items view, but should go back to feeds view
- [ ] "h" and "l" should open and close feeds

## Done
- [x] State should be stored in JSON
- [x] "a" should toggle item read state manually
- [x] "q" quit disabled
- [x] "A" should mark feed as read
- [x] "r" should refresh a single feed
- [x] "Latest" should be displayed in feed item description when present
- [x] Make field test names descriptive
- [x] Table test field names
- [x] Test unloadedFeed
