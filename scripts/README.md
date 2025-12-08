# Organization-Wide Release Drafter Update Script

This directory contains a script to apply the release-drafter configuration changes across all repositories in the `meshery-extensions` organization.

## Script: `update-release-drafter-org-wide.sh`

### Purpose

This script automates the process of updating `.github/release-drafter.yml` configuration files across all active (non-archived, non-forked) repositories in the meshery-extensions organization. It applies the following changes:

1. Adds `no-duplicate-categories: true` to prevent PRs from appearing in multiple release note categories
2. Standardizes label format from singular `label:` to plural `labels:` array format for consistency

### Prerequisites

- **GitHub CLI (`gh`)**: The script requires the GitHub CLI to be installed and authenticated
  - Install: https://cli.github.com/
  - Authenticate: `gh auth login`
- **Bash**: Unix-like shell environment
- **Git**: For repository operations
- **Permissions**: You must have write access to repositories in the meshery-extensions organization

### Installation

1. Install GitHub CLI if not already installed:
   ```bash
   # macOS
   brew install gh
   
   # Linux (Debian/Ubuntu)
   sudo apt install gh
   
   # Windows
   winget install --id GitHub.cli
   ```

2. Authenticate with GitHub:
   ```bash
   gh auth login
   ```

3. Make the script executable (if not already):
   ```bash
   chmod +x scripts/update-release-drafter-org-wide.sh
   ```

### Usage

Run the script from the repository root:

```bash
./scripts/update-release-drafter-org-wide.sh
```

### What the Script Does

1. **Fetches all repositories** in the meshery-extensions organization
2. **Filters** to only process active repositories (excludes archived and forked repos)
3. For each repository:
   - Clones the repository (shallow clone for efficiency)
   - Checks if `.github/release-drafter.yml` exists
   - Checks if the configuration already has `no-duplicate-categories`
   - Creates a new branch (`update-release-drafter-config`)
   - Applies the configuration changes:
     - Adds `no-duplicate-categories: true` after the `tag-template` line
     - Converts singular `label:` to plural `labels:` array format
   - Commits the changes with a detailed message
   - Pushes the branch to the remote
   - Creates a pull request with a comprehensive description
4. **Displays a summary** of processed repositories, PRs created, skipped repos, and any errors

### Output

The script provides colored output showing:
- 🟢 **Green**: Successful operations
- 🟡 **Yellow**: Repositories being processed
- 🔴 **Red**: Errors
- ⊘ **Gray**: Skipped repositories

Example output:
```
Processing: meshery-extensions/meshery-istio
  ✓ Successfully created PR

Summary:
  Total repositories processed: 15
  PRs created: 12
  Skipped: 2
  Errors: 1
```

### Skipped Repositories

The script will skip repositories that:
- Don't have a `.github/release-drafter.yml` file
- Already have `no-duplicate-categories` configured
- Show no changes after the transformation is applied

### Error Handling

The script includes comprehensive error handling and will continue processing other repositories if one fails. Common issues:
- Permission errors (ensure you have write access)
- Network connectivity issues
- Repository already has a PR with the same branch name

### Post-Execution Steps

After the script completes:

1. Review the pull requests created in each repository
2. Ensure the changes look correct
3. Approve and merge the PRs
4. The changes will take effect on the next release draft in each repository

### Configuration Changes Applied

#### Before:
```yaml
name-template: 'Project v$NEXT_PATCH_VERSION'
tag-template: 'v$NEXT_PATCH_VERSION'
categories:
  - title: 📖 Documentation
    label: area/docs
```

#### After:
```yaml
name-template: 'Project v$NEXT_PATCH_VERSION'
tag-template: 'v$NEXT_PATCH_VERSION'
no-duplicate-categories: true
categories:
  - title: 📖 Documentation
    labels:
      - area/docs
```

### Benefits

- **Prevents duplicate entries**: PRs with multiple matching labels now appear only once
- **Maintains categorization**: PRs still get categorized, just in the first matching category
- **Improved consistency**: All repositories use the same configuration pattern
- **Better release notes**: Cleaner, more readable release notes without duplicates

### Troubleshooting

**Issue**: "GitHub CLI (gh) is not installed"
- **Solution**: Install GitHub CLI from https://cli.github.com/

**Issue**: "Not authenticated with GitHub CLI"
- **Solution**: Run `gh auth login` and follow the prompts

**Issue**: "Failed to clone repository"
- **Solution**: Ensure you have read access to the organization's repositories

**Issue**: "Failed to push branch"
- **Solution**: Ensure you have write access to the repositories. Check if a branch with the same name already exists.

**Issue**: "Failed to create PR"
- **Solution**: Check if a PR already exists for this branch. The branch will be pushed successfully but PR creation will fail if one already exists.

### Customization

You can modify the script variables at the top:
- `ORG`: Change the organization name
- `BRANCH_NAME`: Use a different branch name
- `COMMIT_MESSAGE`: Customize the commit message

### Safety Features

- Uses shallow clones (`--depth=1`) for efficiency
- Creates temporary working directory that's cleaned up automatically
- Skips repositories that already have the configuration
- Validates changes before committing
- Provides detailed feedback for each operation
- Continues processing even if individual repositories fail

## Related Documentation

- [Release Drafter Action](https://github.com/release-drafter/release-drafter)
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- [Meshery Extensions Organization](https://github.com/meshery-extensions)

## Support

For issues or questions about this script, please refer to the PR where this script was introduced or open an issue in the repository.
