package config

import (
	"bufio"
	"strings"

	"errors"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	Feeds []string
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ohnurr"), nil
}

func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "feeds"), nil
}

// reads the configuration from disk
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// return empty config if file not found
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Feeds: []string{}}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	c := &Config{
		Feeds: []string{},
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())

		// TODO: check line is a URL

		if l != "" {
			c.Feeds = append(c.Feeds, l)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Save() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// create dir if not exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// will overwrite if file exists
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	writer := bufio.NewWriter(f)
	for _, feed := range c.Feeds {
		_, err = writer.WriteString(feed + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (c *Config) AddFeed(url string) error {
	if slices.Contains(c.Feeds, url) {
		return nil
	}
	c.Feeds = append(c.Feeds, url)
	return nil
}

func (c *Config) RemoveFeed(url string) error {
	for i, feed := range c.Feeds {
		if feed == url {
			c.Feeds = append(c.Feeds[:i], c.Feeds[i+1:]...)
			return nil
		}
	}
	return errors.New("feed not found")
}
