# Contributing to Ohnurr

Thank you for considering contributing to Ohnurr! This document provides guidelines for contributing to the project.

## Development Setup

1. Ensure you have Go 1.21+ installed
2. Clone the repository
3. Run `go mod download` to install dependencies
4. Run `go test ./...` to ensure tests pass

## Commit Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/) for automated changelog generation and semantic versioning.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature (triggers minor version bump)
- **fix**: A bug fix (triggers patch version bump)
- **perf**: Performance improvement
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **style**: Changes that don't affect code meaning (formatting, etc)
- **test**: Adding or updating tests
- **docs**: Documentation only changes
- **chore**: Changes to build process or auxiliary tools
- **ci**: Changes to CI configuration files and scripts

### Examples

```bash
feat: add search functionality to article list
fix: resolve panic when parsing malformed RSS feeds
perf: optimize feed refresh performance
docs: update README with installation instructions
```

### Breaking Changes

Breaking changes should be indicated by a `!` after the type/scope:

```bash
feat!: change config file format to YAML
```

Or include `BREAKING CHANGE:` in the footer:

```bash
feat: update feed parser

BREAKING CHANGE: RSS 1.0 feeds are no longer supported
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Make your changes following the commit conventions
4. Ensure tests pass (`go test ./...`)
5. Ensure code is properly formatted (`go fmt ./...`)
6. Push to your fork and submit a pull request

## Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format code
- Run `go vet` to catch common mistakes
- Keep functions focused and reasonably sized
- Add tests for new functionality

## Testing

- Write tests for new features
- Ensure all tests pass before submitting PR
- Aim for good test coverage

## Reporting Bugs

When reporting bugs, please include:
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

## Suggesting Features

Feature suggestions are welcome! Please:
- Check if the feature has already been suggested
- Clearly describe the use case
- Explain why it would be valuable to users

Thank you for contributing!
