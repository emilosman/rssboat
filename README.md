# rssboat

- **rssboat** is a performant, terminal-based RSS reader written in Go, inspired by [Newsboat](https://github.com/newsboat/newsboat)
- It demonstrates idiomatic Go usage, concurrency, YAML-based configuration, and a TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Features
- Vim-style navigation for feeds and articles  
- Keyboard shortcuts inspired by Newsboat (`?` for help)  
- Asynchronous feed fetching using Go routines and channels  
- Read/unread tracking for feed items  
- Configurable via a simple YAML file  

## Installation
```bash
git clone https://github.com/emilosman/rssboat.git
cd rssboat
go install ./cmd/rssboat
```

## Usage
- Use arrow keys or Vim-style shortcuts to navigate
- Press `?` for the full help menu
- Edit the feed list with your preferred editor (vi by default)
- Mark feeds or items as read/unread

## Configuration (MacOS)
- Config file: `~/Library/Application\ Support/rssboat/urls.yaml`
- Cache file: `~/Library/Caches/rssboat/data.json`

Example urls.yaml:
```
Tech:
  - https://example.com/tech.rss
  - https://example.com/golang.rss
News:
  - https://example.com/worldnews.rss
```

## Development
- This is a hobby project, exploring Go and terminal UI development
- See the [TODO list](./docs/todo.md) for planned features and improvements

## License
- Licensed under [GPLv3](./LICENSE)
