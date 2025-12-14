Ohnurr uses **release-please** for fully automated releases. You just push commits, and everything else happens automatically!

## How It Works

1. **You push commits** using conventional commit format to `main`
2. **release-please automatically creates/updates a Release PR** with:
   - Auto-generated CHANGELOG.md
   - Version bump (based on commit types)
   - All changes since last release
3. **You review and merge the Release PR**
4. **GitHub Actions automatically**:
   - Creates a git tag
   - Runs all tests
   - Builds binaries for all platforms
   - Publishes to GitHub Releases

## Types
- **feat**: New feature (→ shows in changelog, bumps minor)
- **fix**: Bug fix (→ shows in changelog, bumps patch)
- **perf**: Performance improvement (→ shows in changelog)
- **refactor**: Code refactor (→ shows in changelog)
- **docs**: Documentation changes (→ shows in changelog)
- **test**: Test changes (→ hidden from changelog)
- **chore**: Maintenance tasks (→ hidden from changelog)
- **ci**: CI/CD changes (→ hidden from changelog)

## Breaking Changes
```bash
# Option 1: Use ! after type
feat!: completely redesign config format

# Option 2: Use BREAKING CHANGE in footer
feat: add new config system

BREAKING CHANGE: Config format changed from JSON to YAML
```

## CI/CD Workflows

### 1. Continuous Integration (`.github/workflows/ci.yml`)
- **Triggers**: Every push to `main`, every PR
- **Runs**: Tests on Linux/macOS/Windows, linting, build verification

### 2. Release Please (`.github/workflows/release-please.yml`)
- **Triggers**: Every push to `main`
- **Creates**: Release PRs with auto-generated changelogs
- **Publishes**: Release when Release PR is merged

### 3. Manual Release (`.github/workflows/release.yml`)
- **Triggers**: Manual tag push, or workflow dispatch
- **Fallback**: For emergency releases

### Skip a release?
- Close the Release PR without merging
- Commits will accumulate in the next Release PR

### Manual release:
```bash
git tag -a v1.0.1 -m "Hotfix: critical bug"
git push origin v1.0.1
```


## Quick Commands

```bash
# test your changes
go test ./...

# lint code
golangci-lint run

# build locally
go build

# format code
go fmt ./...

# update dependencies
go get -u ./...
go mod tidy
```
