# rssboat
- RSS reader inspired by newsboat

## Todo
- [ ] Unread items should have different color
- [ ] Feeds should load async on UpdateAll()
- [ ] Should send update command after go func finishes
- [ ] Toggle read state on item open
- [ ] "a" should toggle read state manually
- [ ] "A" should mark feed as read
- [ ] Feed items should have a view that displays content
- [ ] feeds.yaml should not be relative path
- [ ] State should be stored in JSON
- [ ] "o" should not trigger feed open when filtering

## Refactor
- [ ] Refactor feed/item open error messages and funcs

## Future
- [ ] Opening feed items should work on all operating systems with default browser

## Maybe
- [ ] "q" should not quit when in items view, but should go back to feeds view
- [ ] "h" and "l" should open and close feeds

## Done
- [x] "r" should refresh a single feed
- [x] "Latest" should be displayed in feed item description when present
- [x] Make field test names descriptive
- [x] Table test field names
- [x] Test unloadedFeed
