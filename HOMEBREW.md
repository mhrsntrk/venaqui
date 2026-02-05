# Homebrew Installation Guide

This guide explains how to set up Venaqui as a Homebrew formula for easy installation.

## Option 1: Using a Homebrew Tap (Recommended)

### Step 1: Create a Tap Repository

Create a new repository named `homebrew-venaqui` (or `homebrew-tap` if you want multiple formulas):

```bash
# Create a new repository on GitHub
# Repository name: homebrew-venaqui
# Make it public
```

### Step 2: Add the Formula

Copy the formula file to your tap repository:

```bash
# In your homebrew-venaqui repository
mkdir -p Formula
cp /path/to/venaqui/Formula/venaqui-simple.rb Formula/venaqui.rb
```

Update the formula with your GitHub repository details:

```ruby
class Venaqui < Formula
  desc "CLI tool with TUI for downloading files via Real-Debrid and aria2"
  homepage "https://github.com/mhrsntrk/venaqui"
  url "https://github.com/mhrsntrk/venaqui/archive/v0.1.0.tar.gz"
  sha256 "YOUR_SHA256_HASH"
  license "MIT"
  head "https://github.com/mhrsntrk/venaqui.git", branch: "main"

  depends_on "go" => :build
  depends_on "aria2"

  def install
    system "go", "build", "-o", bin/"venaqui", "./cmd/venaqui"
  end

  test do
    system "#{bin}/venaqui", "version"
  end
end
```

### Step 3: Calculate SHA256

When you create a release, calculate the SHA256:

```bash
# Download the source tarball
wget https://github.com/mhrsntrk/venaqui/archive/v0.1.0.tar.gz

# Calculate SHA256
shasum -a 256 v0.1.0.tar.gz
```

Update the formula with the SHA256 hash.

### Step 4: Install via Tap

Users can now install with:

```bash
brew tap mhrsntrk/venaqui
brew install venaqui
```

## Option 2: Using GoReleaser (Automated)

If you use GoReleaser for releases, it can automatically generate and update the Homebrew formula.

### Step 1: Configure GoReleaser

The `.goreleaser.yml` file is already configured with Homebrew support. When you create a release:

```bash
# Tag your release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# GoReleaser will create the release and formula
goreleaser release
```

### Step 2: Create Tap Repository

Create `homebrew-venaqui` repository and add the generated formula from GoReleaser output.

### Step 3: Update Formula

GoReleaser will output the formula. Copy it to your tap repository.

## Option 3: Local Installation (Development)

For local testing:

```bash
# Install from local formula
brew install --build-from-source Formula/venaqui-simple.rb

# Or create a local tap
brew tap-new mhrsntrk/venaqui
cp Formula/venaqui-simple.rb $(brew --repository mhrsntrk/venaqui)/Formula/venaqui.rb
brew install venaqui
```

## Updating the Formula

When releasing a new version:

1. Update the version in the formula
2. Update the URL to point to the new release
3. Calculate and update the SHA256 hash
4. Commit and push to your tap repository

## Testing

Test your formula:

```bash
# Install
brew install mhrsntrk/venaqui/venaqui

# Test
venaqui version

# Uninstall
brew uninstall venaqui
```

## Troubleshooting

### Formula not found

Make sure you've tapped the repository:
```bash
brew tap mhrsntrk/venaqui
```

### SHA256 mismatch

Recalculate the SHA256 for your release tarball and update the formula.

### Build errors

Ensure all dependencies are listed in the formula:
- `go` (for building)
- `aria2` (runtime dependency)
