package main

import (
	"fmt"
	"os"

	"ohnurr/config"
	"ohnurr/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		// no args so launch TUI
		launchTUI()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ohnurr add <url>")
			os.Exit(1)
		}
		addFeed(os.Args[2])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ohnurr remove <url>")
			os.Exit(1)
		}
		removeFeed(os.Args[2])
	case "list":
		listFeeds()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func addFeed(url string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.AddFeed(url); err != nil {
		fmt.Printf("Error adding feed: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Save(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added feed: %s\n", url)
}

func removeFeed(url string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.RemoveFeed(url); err != nil {
		fmt.Printf("Error removing feed: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Save(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed feed: %s\n", url)
}

func listFeeds() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Feeds) == 0 {
		fmt.Println("No feeds configured")
		return
	}

	fmt.Println("Configured feeds:")
	for i, feed := range cfg.Feeds {
		fmt.Printf("%d. %s\n", i+1, feed)
	}
}

func launchTUI() {
	c, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(c.Feeds) == 0 {
		fmt.Println("No feeds configured. Add some feeds first:")
		fmt.Println("  ohnurr add <url>")
		os.Exit(1)
	}

	state, err := config.LoadState()
	if err != nil {
		fmt.Printf("Error loading state: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(c, state)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("ohnurr - Terminal RSS Reader")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ohnurr              Launch interactive TUI")
	fmt.Println("  ohnurr add <url>    Add RSS feed")
	fmt.Println("  ohnurr remove <url> Remove RSS feed")
	fmt.Println("  ohnurr list         List all feeds")
	fmt.Println("  ohnurr help         Show this help message")
}
