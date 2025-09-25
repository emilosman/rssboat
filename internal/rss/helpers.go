package rss

import (
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func LoadList(filesystem fs.FS) (*List, error) {
	l := List{
		FeedIndex: make(map[string]*RssFeed),
	}

	err := l.CreateFeedsFromYaml(filesystem, "urls.yaml")
	if err != nil {
		return &l, err
	}

	cacheFilePath, err := CacheFilePath()
	if err != nil {
		return &l, err
	}

	f, err := os.Open(cacheFilePath)
	if err != nil {
		return &l, err
	}
	defer f.Close()

	err = l.Restore(f)
	if err != nil {
		return &l, err
	}

	return &l, nil
}

func CacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "rssboat")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "data.json"), nil
}

func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "rssboat")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	configFile := filepath.Join(appDir, "urls.yaml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
		if err != nil {
			return "", err
		}
		defer f.Close()
	}

	return appDir, nil
}

func Clean(input string) string {
	p := bluemonday.StrictPolicy()
	clean := p.Sanitize(input)
	decoded := html.UnescapeString(clean)
	return normalizeSpaces(decoded)
}

func normalizeSpaces(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	return strings.Join(strings.Fields(s), " ")
}
