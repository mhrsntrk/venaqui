# GitHub Repository Setup

This guide will help you push Venaqui to GitHub and set up the repository.

## Step 1: Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `venaqui`
3. Description: "CLI tool with TUI for downloading files via Real-Debrid and aria2"
4. Visibility: Public (or Private, your choice)
5. **DO NOT** initialize with README, .gitignore, or license (we already have these)
6. Click "Create repository"

## Step 2: Add Remote and Push

After creating the repository, run these commands:

```bash
cd /Users/mhrsntrk/Developer/venaqui

# Add the remote repository
git remote add origin https://github.com/mhrsntrk/venaqui.git

# Push to GitHub
git push -u origin main
```

If you're using SSH instead of HTTPS:

```bash
git remote add origin git@github.com:mhrsntrk/venaqui.git
git push -u origin main
```

## Step 3: Verify

Visit https://github.com/mhrsntrk/venaqui to verify everything is pushed correctly.

## Step 4: Set Up GitHub Actions (Optional)

The repository includes a GitHub Actions workflow for automated releases. To enable it:

1. Go to repository Settings → Actions → General
2. Enable "Workflow permissions" → Read and write permissions
3. When you create your first release tag, the workflow will automatically run

## Step 5: Create First Release

When ready to create your first release:

```bash
# Tag the release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# Or use GitHub web interface:
# Go to Releases → Draft a new release → Create new tag v0.1.0
```

## Step 6: Set Up Homebrew Tap (Later)

After your first release, create a separate repository for the Homebrew tap:

1. Create new repository: `homebrew-venaqui`
2. Copy the formula from `Formula/venaqui-simple.rb`
3. Update the formula with the correct URL and SHA256
4. See `HOMEBREW.md` for detailed instructions

## Troubleshooting

### Authentication Issues

If you get authentication errors:

```bash
# Use GitHub CLI (if installed)
gh auth login

# Or set up SSH keys
# See: https://docs.github.com/en/authentication/connecting-to-github-with-ssh
```

### Push Rejected

If push is rejected, you might need to pull first:

```bash
git pull origin main --allow-unrelated-histories
```
