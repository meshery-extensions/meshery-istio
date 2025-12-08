#!/bin/bash

# Script to update release-drafter.yml configuration across all repositories in meshery-extensions organization
# This script adds 'no-duplicate-categories: true' and standardizes label format for consistency

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed. Please install it from https://cli.github.com/${NC}"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${RED}Error: Not authenticated with GitHub CLI. Run 'gh auth login' first.${NC}"
    exit 1
fi

# Configuration
ORG="meshery-extensions"
BRANCH_NAME="update-release-drafter-config"
COMMIT_MESSAGE="Add no-duplicate-categories to prevent duplicate PR entries in release notes"

echo -e "${GREEN}Starting release-drafter configuration update for ${ORG} organization...${NC}\n"

# Get all repositories in the organization
echo "Fetching repositories from ${ORG}..."
repos=$(gh repo list "$ORG" --limit 1000 --json name,isArchived,isFork --jq '.[] | select(.isArchived == false and .isFork == false) | .name')

if [ -z "$repos" ]; then
    echo -e "${RED}No repositories found or unable to access organization.${NC}"
    exit 1
fi

repo_count=$(echo "$repos" | wc -l)
echo -e "${GREEN}Found ${repo_count} active repositories in ${ORG}${NC}\n"

# Create temporary directory for work
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

processed=0
updated=0
skipped=0
errors=0

# Process each repository
while IFS= read -r repo; do
    echo -e "\n${YELLOW}Processing: ${ORG}/${repo}${NC}"
    ((processed++))
    
    cd "$TEMP_DIR"
    
    # Clone the repository
    if ! gh repo clone "${ORG}/${repo}" "${repo}" -- --depth=1 2>/dev/null; then
        echo -e "${RED}  ✗ Failed to clone repository${NC}"
        ((errors++))
        continue
    fi
    
    cd "$repo"
    
    # Get the default branch for this repository using gh CLI
    default_branch=$(gh repo view "${ORG}/${repo}" --json defaultBranchRef --jq '.defaultBranchRef.name' 2>/dev/null)
    if [ -z "$default_branch" ]; then
        # Fallback: try to detect from git
        default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
        if [ -z "$default_branch" ]; then
            # Last resort: check if main exists, otherwise use master
            if git show-ref --verify --quiet refs/remotes/origin/main; then
                default_branch="main"
            else
                default_branch="master"
            fi
        fi
    fi
    
    # Check if release-drafter.yml exists
    if [ ! -f ".github/release-drafter.yml" ]; then
        echo -e "  ⊘ Skipped: No .github/release-drafter.yml file found"
        ((skipped++))
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    fi
    
    # Check if already has no-duplicate-categories
    if grep -q "no-duplicate-categories:" ".github/release-drafter.yml"; then
        echo -e "  ⊘ Skipped: Already has no-duplicate-categories configured"
        ((skipped++))
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    fi
    
    # Create a new branch
    git checkout -b "$BRANCH_NAME" 2>/dev/null || {
        echo -e "${RED}  ✗ Failed to create branch${NC}"
        ((errors++))
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    }
    
    # Backup original file
    cp ".github/release-drafter.yml" ".github/release-drafter.yml.bak"
    
    # Check if tag-template exists in the file
    if ! grep -q "^tag-template:" ".github/release-drafter.yml.bak"; then
        echo -e "${RED}  ✗ No 'tag-template:' found in release-drafter.yml${NC}"
        ((errors++))
        rm ".github/release-drafter.yml.bak"
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    fi
    
    # Add no-duplicate-categories after tag-template line
    if ! awk '/^tag-template:/ {print; print "no-duplicate-categories: true"; next} 1' \
        ".github/release-drafter.yml.bak" > ".github/release-drafter.yml"; then
        echo -e "${RED}  ✗ Failed to update configuration${NC}"
        ((errors++))
        rm ".github/release-drafter.yml.bak"
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    fi
    
    # Verify that no-duplicate-categories was added
    if ! grep -q "^no-duplicate-categories:" ".github/release-drafter.yml"; then
        echo -e "${RED}  ✗ Failed to insert no-duplicate-categories line${NC}"
        ((errors++))
        rm ".github/release-drafter.yml.bak"
        cd "$TEMP_DIR"
        rm -rf "$repo"
        continue
    fi
    
    # Fix label vs labels inconsistency (singular to plural array format)
    # This handles cases where 'label:' is used instead of 'labels:'
    # Detects indentation dynamically (2 or 4 spaces)
    # Using a more portable approach with a temporary file
    while IFS= read -r line; do
        # Check for 2-space indented label
        if echo "$line" | grep -q "^  label: "; then
            label_value=$(echo "$line" | sed 's/^  label: //')
            printf "  labels:\n    - %s\n" "$label_value"
        # Check for 4-space indented label
        elif echo "$line" | grep -q "^    label: "; then
            label_value=$(echo "$line" | sed 's/^    label: //')
            printf "    labels:\n      - %s\n" "$label_value"
        else
            printf "%s\n" "$line"
        fi
    done < ".github/release-drafter.yml" > ".github/release-drafter.yml.tmp"
    mv ".github/release-drafter.yml.tmp" ".github/release-drafter.yml"
    
    # Clean up backup
    rm ".github/release-drafter.yml.bak"
    
    # Check if there are changes
    if ! git diff --quiet ".github/release-drafter.yml"; then
        # Commit changes
        git add ".github/release-drafter.yml"
        git commit -m "$COMMIT_MESSAGE" \
            -m "- Added 'no-duplicate-categories: true' to prevent PRs from appearing in multiple categories" \
            -m "- Standardized label format to use 'labels:' array format for consistency" \
            -m "" \
            -m "When a PR has multiple labels matching different categories, it will now appear only in the first matching category." 2>/dev/null || {
            echo -e "${RED}  ✗ Failed to commit changes${NC}"
            ((errors++))
            cd "$TEMP_DIR"
            rm -rf "$repo"
            continue
        }
        
        # Push branch
        if git push origin "$BRANCH_NAME" 2>/dev/null; then
            # Create pull request using a here-document for better readability
            pr_body=$(cat <<'EOF_PR_BODY'
**Description**

This PR adds the `no-duplicate-categories: true` configuration option to the release-drafter configuration file. This ensures that each merged pull request is only included at most once in the release notes, even if it has multiple labels that match different categories.

**Changes:**
- Added `no-duplicate-categories: true` to `.github/release-drafter.yml`
  - PRs now appear only in the first matching category based on defined order
- Standardized any singular `label:` format to use `labels:` array format for consistency

**Notes for Reviewers**

Configuration-only change. When merged, the next release draft will include each PR at most once.

---
This change is being applied across all repositories in the organization.
EOF_PR_BODY
)
            
            if gh pr create \
                --title "[release-drafter] Prevent duplicate PR entries in release notes" \
                --body "$pr_body" \
                --base "$default_branch" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Successfully created PR${NC}"
                ((updated++))
            else
                echo -e "${RED}  ✗ Failed to create PR (branch pushed successfully)${NC}"
                ((errors++))
            fi
        else
            echo -e "${RED}  ✗ Failed to push branch${NC}"
            ((errors++))
        fi
    else
        echo -e "  ⊘ Skipped: No changes detected after transformation"
        ((skipped++))
    fi
    
    # Cleanup
    cd "$TEMP_DIR"
    rm -rf "$repo"
    
done <<< "$repos"

# Summary
echo -e "\n${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}Summary:${NC}"
echo -e "  Total repositories processed: ${processed}"
echo -e "  ${GREEN}PRs created: ${updated}${NC}"
echo -e "  ${YELLOW}Skipped: ${skipped}${NC}"
echo -e "  ${RED}Errors: ${errors}${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}\n"

if [ $updated -gt 0 ]; then
    echo -e "${GREEN}Successfully created ${updated} pull request(s) across the organization!${NC}"
    echo -e "\nNext steps:"
    echo -e "  1. Review the PRs created in each repository"
    echo -e "  2. Merge the PRs after approval"
    echo -e "  3. The changes will take effect on the next release draft"
fi

exit 0
