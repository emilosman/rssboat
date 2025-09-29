# Todo

## MVP
- [ ] [Help keys extend](https://chatgpt.com/c/68c1ad14-5c9c-8331-bad6-ce4f7c1f52c8)

## Todo
- [ ] Don't delete cache on error
- shift+e
  - [ ] Custom Cmd.Msg while editing
  - [ ] Ignore key events when editing in vim (check custom Cmd.Msg)
  - [ ] Default editor instead of vim
- [ ] Single feed update complete message (instead of "All feeds updated")
- [ ] 'n' jump to next unread
- [ ] Updating / Updated message reformat. Show both messages
- [ ] Tab unread indicator tab (green .)
- [ ] Parse URL with standard library to check for errors
- [ ] Content display view
- [ ] 100% test coverage of rss.go
- [ ] 🔴 error feed
- [ ] Clear item list filter on back
- Tabs
  - [ ] "All" tab
  - [ ] Tab change on 1-9
  - [ ] Test long feed list for "h" and "l" handling
- urls.yaml
  - [ ] urls.yaml custom ENV path support
  - [ ] Newsboat urls.txt support - read from ~/.newsboat/urls ? - modal dialog ?
- [ ] Timeout network request 8s
- [ ] Unread items should have different color
- [ ] Unread counter (15/254)
- [ ] Reset selected item index when opening feed, remember it on the feedlist
- [ ] README extend
- [ ] Record demo using Charm's vhs

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
- [ ] E2E tests [teatest](https://github.com/caarlos0/teatest-example/blob/main/main_test.go)

## Ownership
- [ ] Transfer rssboat to picigato org
- [ ] GH sponsors setup finish [link](https://github.com/sponsors/picigato/signup)

## Maybe
- [ ] Preserve tab order from urls.yaml
- [ ] Index number in front of items
- [ ] Confirmation Y/N on major commands (Mark all feeds as read...)
- [ ] "h" and "l" should open and close feeds
- [ ] Sort options (1. newest unread up top, 2. popular...)
- [ ] Remember tab selection on close

## Done
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
