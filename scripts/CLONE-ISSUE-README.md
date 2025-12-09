# Issue Cloner Script

This script clones (replicates) a GitHub issue to multiple repositories within an organization. The cloned issues will have the same title, body, labels, assignees, and milestone as the source issue.

## Script: `clone-issue-org-wide.sh`

### Purpose

This script allows you to replicate a GitHub issue across all active (non-archived, non-forked) repositories in an organization. This is useful when you need to:

- Track the same task or bug across multiple repositories
- Announce organization-wide changes that require action in each repository
- Coordinate updates or migrations across all projects
- Ensure consistency of issues across related repositories

### Features

- Clones issue title, body, labels, assignees, and milestone
- Supports custom target organization (or uses source issue's organization)
- Automatically skips the source repository to avoid duplicates
- Skips repositories with issues disabled
- Adds a reference to the source issue in the cloned issue body
- Interactive confirmation before creating issues
- Detailed progress reporting with color-coded output
- Comprehensive error handling

### Prerequisites

- **GitHub CLI (`gh`)**: The script requires the GitHub CLI to be installed and authenticated
  - Install: https://cli.github.com/
  - Authenticate: `gh auth login`
- **jq**: JSON processor for parsing issue data
  - Install: `sudo apt-get install jq` (Linux) or `brew install jq` (macOS)
- **Bash**: Unix-like shell environment
- **Permissions**: You must have write access to repositories in the target organization

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

2. Install jq:
   ```bash
   # macOS
   brew install jq
   
   # Linux (Debian/Ubuntu)
   sudo apt-get install jq
   ```

3. Authenticate with GitHub:
   ```bash
   gh auth login
   ```

4. Make the script executable (if not already):
   ```bash
   chmod +x scripts/clone-issue-org-wide.sh
   ```

### Usage

**Basic usage (clone to same organization):**
```bash
./scripts/clone-issue-org-wide.sh <issue-url>
```

**Clone to different organization:**
```bash
./scripts/clone-issue-org-wide.sh <issue-url> <target-organization>
```

### Examples

**Example 1: Clone issue within the same organization**
```bash
./scripts/clone-issue-org-wide.sh https://github.com/meshery-extensions/meshery-istio/issues/123
```
This will clone issue #123 from meshery-istio to all other repositories in the meshery-extensions organization.

**Example 2: Clone issue to a different organization**
```bash
./scripts/clone-issue-org-wide.sh https://github.com/meshery-extensions/meshery-istio/issues/123 my-other-org
```
This will clone issue #123 to all repositories in the `my-other-org` organization.

### What the Script Does

1. **Validates input**: Checks that the issue URL is in the correct format
2. **Fetches issue details**: Retrieves title, body, labels, assignees, milestone, and state from the source issue
3. **Fetches repositories**: Gets all active (non-archived, non-forked) repositories in the target organization
4. **Confirms action**: Asks for user confirmation before creating issues
5. For each repository:
   - Skips the source repository (to avoid duplicates)
   - Checks if issues are enabled
   - Creates a new issue with the same attributes as the source
   - Adds a reference note at the bottom linking back to the source issue
6. **Reports summary**: Shows total created, skipped, and errors

### Output

The script provides color-coded output:
- 🔵 **Blue**: Informational messages
- 🟢 **Green**: Successful operations
- 🟡 **Yellow**: Warnings or skipped items
- 🔴 **Red**: Errors

Example output:
```
═══════════════════════════════════════════════════════════
GitHub Issue Cloner
═══════════════════════════════════════════════════════════
Source: meshery-extensions/meshery-istio#123
Target Organization: meshery-extensions

Fetching issue details...
✓ Issue details fetched
  Title: Update documentation for new feature
  State: OPEN
  Labels: kind/docs,priority/high
  Assignees: username
  Milestone: v1.0

Fetching repositories from meshery-extensions...
Found 18 active repositories in meshery-extensions

This will create 18 new issues across the organization.
Do you want to continue? (y/N): y

⊘ Skipped: meshery-extensions/meshery-istio (source repository)
Processing: meshery-extensions/meshery-consul
  ✓ Issue created: https://github.com/meshery-extensions/meshery-consul/issues/45
Processing: meshery-extensions/meshery-linkerd
  ✓ Issue created: https://github.com/meshery-extensions/meshery-linkerd/issues/67
...

═══════════════════════════════════════════════════════════
Summary:
  Total repositories processed: 18
  Issues created: 17
  Skipped: 1
  Errors: 0
═══════════════════════════════════════════════════════════

Successfully created 17 issue(s) across the organization!
```

### Cloned Issue Format

Each cloned issue will have:
- **Same title** as the source issue
- **Same body** as the source issue, with a footer note added:
  ```
  ---
  _This issue was cloned from owner/repo#number_
  ```
- **Same labels** (if they exist in the target repository)
- **Same assignees** (if they have access to the target repository)
- **Same milestone** (if it exists in the target repository)

### Skipped Repositories

The script will skip repositories that:
- Are archived
- Are forks
- Are the source repository (to avoid duplicates)
- Have issues disabled

### Error Handling

The script includes comprehensive error handling and will:
- Continue processing other repositories if one fails
- Display detailed error messages for debugging
- Provide a summary of successes, skips, and errors

Common errors and solutions:

**Issue creation failed**
- **Cause**: Labels, assignees, or milestone don't exist in target repository
- **Solution**: The issue will still be created without those attributes. You can add them manually after creation.

**Permission denied**
- **Cause**: Insufficient permissions in target repository
- **Solution**: Ensure you have write access to the repository

**Issues not enabled**
- **Cause**: Repository has issues feature disabled
- **Solution**: Enable issues in repository settings, or the script will skip it

### Safety Features

- **Interactive confirmation**: Asks before creating issues to prevent accidental mass creation
- **Source repository skip**: Automatically skips the source repository to avoid duplicates
- **Validation**: Checks issue URL format and authentication before proceeding
- **Error isolation**: Failures in one repository don't affect others

### Use Cases

1. **Organization-wide announcements**: Create an issue in all repos about upcoming changes
2. **Security updates**: Track security patches needed across all projects
3. **Migration tasks**: Coordinate dependency updates or API migrations
4. **Compliance requirements**: Ensure all repositories address compliance items
5. **Documentation updates**: Request documentation updates across all projects

### Limitations

- Labels, assignees, and milestones must exist in target repositories (script will create issue without them if they don't exist)
- Cannot copy issue comments (only the original issue body)
- Cannot copy issue reactions or other metadata
- Requires manual intervention if you want to link the cloned issues back to each other

### Tips

- Review the issue content before cloning to ensure it's appropriate for all repositories
- Consider adding organization-specific context to the issue body before cloning
- Use clear labels that exist across all repositories for better categorization
- Add instructions in the issue body specific to what needs to be done in each repository

## Related Scripts

- `update-release-drafter-org-wide.sh`: Updates release-drafter configuration across repositories
- See main `README.md` for other automation scripts

## Support

For issues or questions about this script, please refer to the repository where it was introduced or open an issue in the appropriate repository.
