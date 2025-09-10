# rssboat
- RSS reader inspired by newsboat

## MVP
- [ ] urls.yaml should not be relative path
  - [ ] Should handle urls.yaml not existing

## Todo
- [ ] [Help keys extend](https://chatgpt.com/c/68c1ad14-5c9c-8331-bad6-ce4f7c1f52c8)
- [ ] Parse URL with standard library to check for errors
- Tabs
  - [ ] "All" tab
  - [ ] Tab change on 1-9 or Tab
  - [ ] Test long feed list for "h" and "l" handling
  - [ ] Refresh tab only "ctrl+r"
  - [ ] Unread indicator tab
- [ ] Edit urls.yaml with "shift+e"
- [ ] Single feed refresh status messages fix
- [ ] CategoriesIndex on List
- [ ] Content display view
- [ ] Unread counter (15/254)
- [ ] Unit test tui funcs (BuildFeedList...)
- [ ] Timeout network request 10s
- [ ] Unread items should have different color
- [ ] Should send update command after go func finishes
- [ ] Reset selected item index when opening feed, remember it on the feedlist
- [ ] 100% test coverage of rss.go
- [ ] Record demo using Charm's vhs
- [ ] Fix linebreak in item titles

## Refactor
- [ ] Refactor feed list title usage + custom status message display ?
- [ ] Refactor GetCategory
- [ ] Refactor list updates into helper func
- [ ] Refactor feed/item open error messages and funcs

## Future
- [ ] Public repo
- [ ] GH sponsors
- [ ] Reddit post
- [ ] E2E tests [teatest](https://github.com/caarlos0/teatest-example/blob/main/main_test.go)

## Maybe
- [ ] Confirmation Y/N on major commands
- [ ] "h" and "l" should open and close feeds
- [ ] Sort options (1. newest unread up top, 2. popular...)

## Done
- [x] License
- [x] Sanitize content (titles, items, descriptions) of HTML, JS, newlines...
- [x] Handle update event for async
- [x] Feeds should load async on UpdateAll(), update one-by-one
- [x] Color active tab
- [x] Tabs for categories
  - [example](https://github.com/charmbracelet/bubbletea/blob/28ab4f41b29fef14d900c46a4873a45891a9ee9b/examples/tabs/main.go#L40)
- [x] FeedIndex
- [x] Check YAML file when restoring from JSON
- [x] Messages moved to vars
- [x] DRY cleanup
- [x] Refactor use of pointers
- [x] Increase test coverage
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


## Paths

### MacOS
- cache: ~/Library/Caches/rssboat/data.json
- config: ~/Library/Application\ Support/rssboat/urls.yaml
