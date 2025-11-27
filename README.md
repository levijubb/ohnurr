# Ohnurr

Ohnurr a TUI RSS reader.

> **Note:** This project is very much in development and may have bugs. PRs are welcome.

## Prerequisites

- Go 1.21+ (for `go build`)
- A terminal that supports ANSI colors

## Build & Run
Ensure your the go bin directory is in your $PATH:

unix `$HOME/go/bin`

windows `%USERPROFILE%\\go\\bin`

```bash
go build -o ohnurr
go install
```

## Commands
```bash
ohnurr              # launch TUI
ohnurr add <url>    # add a feed
ohnurr remove <url> # drop a feed
ohnurr list         # show configured feeds
```

`feeds` and `state` files are plain text and live in `~/.config/ohnurr/`.
