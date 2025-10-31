# Todo

## MVP
- [ ] Create default values in urls.yaml
- [ ] README extend
  - [ ] Explain commands in readme: ctrl+ tabs, shift+ all items
- [ ] Test first run flow (lists_test.go)

## Todo
- [ ] Viewport full help heigh fix (vertical join?)
- [ ] 100% test coverage of rss.go
- [ ] Items list height bug in MacOS Terminal
- [ ] Record demo using Charm's vhs
- [ ] Single feed update complete message (instead of "All feeds updated")
- shift+e
  - [ ] Fix ENV support for default editor
- [ ] Long feed list hide/disable "l", "h", "pgdwn", "pgup" display in help
- [ ] Updating / Updated message reformat. Show both messages
- [ ] Parse URL with standard library to check for errors
- urls.yaml
  - [ ] urls.yaml custom ENV path support
  - [ ] Newsboat urls.txt support - read from ~/.newsboat/urls ? - modal dialog ?
- [ ] Timeout network request 8s
- [ ] Unread counter (15/254)
- [ ] Linux test

## Database
- [ ] Use database instead of JSON only

## Refactor
- Categories
  - [ ] Refactor GetCategory
  - [ ] CategoriesIndex on List
- [ ] Refactor m.selectedFeed vs m.lf.SelectedItem() usage (handleMarkFeedRead)
- shift+e
  - [ ] Refactor handleEdit to not have tui depend on rss package
  - [ ] Refactor handleEdit + initialModel DRY
  - [ ] Refactor initial list build (issues with tab rebuild, data.json items remain after items removed)
- [ ] Config paths func should not be in the rss package
- [ ] Table driven tests

## Future
- [ ] Reddit post
  - [ ] r/rss
  - [ ] r/newsboat
  - [ ] r/golang
- [ ] E2E tests [teatest](https://github.com/caarlos0/teatest-example/blob/main/main_test.go)
- [ ] @latest

## Ownership
- [ ] Transfer rssboat to picigato org
- [ ] GH sponsors setup finish [link](https://github.com/sponsors/picigato/signup)

## Maybe
- [ ] "All" tab (tab 0 ?)
- [ ] Capslock warning
- [ ] Select next item after "a" mark as read toggle
- [ ] Preserve tab order from urls.yaml
- [ ] Index number in front of items
- [ ] Jump to line (e.g.: `:2`) vim style
- [ ] Confirmation Y/N on major commands (Mark all feeds as read...)
- [ ] "h" and "l" should open and close feeds
- [ ] Sort options (1. newest unread up top, 2. popular...)
- [ ] Remember tab selection on close

## Done
- [x] Item search Title + Content
- [x] UpdateStatusMessage auto-clear
- [x] Clear status after N seconds
- [x] Custom list title + status message
- [x] Title + status in viewport
- [x] Viewport help short and full list different
- [x] Toggle read viewport handler
- [x] Marked as read / unread message instead of "Read state toggled"
- [x] Item delete from urls.yaml does not remove from list?
- [x] Fix feeds getting marked unread after edit
- [x] Save command check usage
- [x] Item preview refinements
- [x] Viewport help modal
- [x] Reset selected item index when opening feed, remember it on the feedlist
- [x] "p" previous unread feed
- [x] "p" previous unread item
- [x] 'n' jump to next unread
- [x] Unread items should have different color
- [x] Next unread feed
- [x] Next article in viewport mode "l", "right"
- [x] Prev article in viewport mode "h", "left"
- [x] Next unread item
- [x] Can open cached errored feeds
- [x] Tab change on 1-9
- [x] Red error feed
- [x] Clear item list filter on back
- [x] Display "+" instead of 🟢 in terminals that don't support it
- [x] Content display view
- [x] First use message, empty list, show "shift+e" prompt
- [x] Fix browser open
- [x] Fresh install test
- [x] Windows test
- [x] Test long feed list for "h" and "l" handling
- [x] Help keys extend
- [x] Mark all tab items read "ctrl+a"
- [x] Update tab only "ctrl+r"
- [x] Single feed update refactor, status messages fix
- [x] Rebuild tabs after shift+e
- [x] Break up rss_test.go
- [x] rss.go break up into smaller files
- [x] Public repo
- [x] README.md write
- [x] Categories order alphabetical or as in yaml
- [x] Edit urls.yaml with "shift+e"
- [x] "Latest" should return first unread if present
- [x] Refactor feed/item open error messages and funcs
- [x] Refactor list updates into helper func
- [x] Refactor feed list title usage + custom status message display
- [x] Should handle urls.yaml not existing
- [x] Unit test tui funcs (BuildFeedList...)
- [x] Should send update command after go func finishes
- [x] Fix linebreak in item titles
- [x] urls.yaml should be absolute path
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
