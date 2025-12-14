# Ohnurr

[![CI](https://github.com/levijubb/ohnurr/actions/workflows/ci.yml/badge.svg)](https://github.com/levijubb/ohnurr/actions/workflows/ci.yml)
[![Release](https://github.com/levijubb/ohnurr/actions/workflows/release.yml/badge.svg)](https://github.com/levijubb/ohnurr/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/levijubb/ohnurr)](https://goreportcard.com/report/github.com/levijubb/ohnurr)
[![License](https://img.shields.io/github/license/levijubb/ohnurr)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/levijubb/ohnurr)](https://github.com/levijubb/ohnurr/releases/latest)

Ohnurr is a modern TUI (Terminal User Interface) RSS reader built with Go.

> **Note:** This project is very much in development and may have bugs. Feel free to open PRs.

## Installation

### Using Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/levijubb/ohnurr/releases/latest).

**macOS/Linux:**
```bash
# Extract and move to PATH
tar -xzf ohnurr_*_*.tar.gz
sudo mv ohnurr /usr/local/bin/
```

**Windows:**
Extract the zip file and add `ohnurr.exe` to your PATH.

### Using Go Install

```bash
go install github.com/levijubb/ohnurr@latest
```

### Build from Source

**Prerequisites:**
- Go 1.21+ (for `go build`)
- A terminal that supports ANSI colors

Ensure the go bin directory is in your $PATH:
- Unix: `$HOME/go/bin`
- Windows: `%USERPROFILE%\go\bin`

```bash
git clone https://github.com/levijubb/ohnurr.git
cd ohnurr
go build -o ohnurr
go install
```

## Usage

```bash
ohnurr                 # Launch interactive TUI
ohnurr add <url>       # Add RSS feed
ohnurr remove <url>    # Remove RSS feed
ohnurr list            # List all feeds
ohnurr version         # Show version information
ohnurr help            # Show help message
```

Configuration files (`feeds` and `state`) are stored in `~/.config/ohnurr/`.
