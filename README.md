# Ohnurr

Ohnurr a TUI RSS reader.

## Prerequisites

- Go 1.21+ (for `go build`)
- A terminal that supports ANSI colors

## Build & Run

```bash
go build -o ohnurr
./ohnurr            # launch TUI
ohnurr add <url>    # add a feed
ohnurr remove <url> # drop a feed
ohnurr list         # show configured feeds
```

`feeds` and `state` files are plain text and live in `~/.config/ohnurr/`.

