package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type State struct {
	ReadArticles map[string]bool // Key is article GUID or link
}

func GetStatePath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state"), nil
}

// reads the state from disk
func LoadState() (*State, error) {
	path, err := GetStatePath()
	if err != nil {
		return nil, err
	}

	// if state doesn't exist, return empty state
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &State{ReadArticles: make(map[string]bool)}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	state := &State{
		ReadArticles: make(map[string]bool),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		if l != "" {
			state.ReadArticles[l] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return state, nil
}

// save writes the state to disk
func (s *State) Save() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := GetStatePath()
	if err != nil {
		return err
	}

	// create/overwrite the file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for articleID := range s.ReadArticles {
		_, err = writer.WriteString(articleID + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (s *State) MarkAsRead(articleID string) {
	s.ReadArticles[articleID] = true
}

func (s *State) IsRead(articleID string) bool {
	return s.ReadArticles[articleID]
}
