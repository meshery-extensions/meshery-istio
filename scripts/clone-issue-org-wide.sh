#!/bin/bash

# Script to clone a GitHub issue to multiple repositories in an organization
# This creates a replica of the source issue with the same title, body, labels, assignees, and milestone

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Usage function
usage() {
    echo "Usage: $0 <issue-url> [target-organization]"
    echo ""
    echo "Arguments:"
    echo "  issue-url              Full URL to the source GitHub issue (e.g., https://github.com/org/repo/issues/123)"
    echo "  target-organization    Optional. Organization to clone issue to. If not specified, uses the source issue's organization."
    echo ""
    echo "Examples:"
    echo "  $0 https://github.com/meshery-extensions/meshery-istio/issues/123"
    echo "  $0 https://github.com/meshery-extensions/meshery-istio/issues/123 my-org"
    exit 1
}

# Parse arguments
if [ $# -lt 1 ]; then
    usage
fi

ISSUE_URL="$1"
TARGET_ORG="${2:-}"

# Parse issue URL to extract owner, repo, and issue number
if [[ ! "$ISSUE_URL" =~ ^https://github\.com/([^/]+)/([^/]+)/issues/([0-9]+)$ ]]; then
    echo -e "${RED}Error: Invalid issue URL format. Expected: https://github.com/owner/repo/issues/number${NC}"
    exit 1
fi

SOURCE_OWNER="${BASH_REMATCH[1]}"
SOURCE_REPO="${BASH_REMATCH[2]}"
ISSUE_NUMBER="${BASH_REMATCH[3]}"

# If target org not specified, use source owner
if [ -z "$TARGET_ORG" ]; then
    TARGET_ORG="$SOURCE_OWNER"
fi

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}GitHub Issue Cloner${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "Source: ${SOURCE_OWNER}/${SOURCE_REPO}#${ISSUE_NUMBER}"
echo -e "Target Organization: ${TARGET_ORG}"
echo ""

# Fetch issue details
echo "Fetching issue details..."
issue_data=$(gh issue view "$ISSUE_NUMBER" --repo "${SOURCE_OWNER}/${SOURCE_REPO}" --json title,body,labels,assignees,milestone,state 2>&1)
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to fetch issue details${NC}"
    echo -e "${RED}${issue_data}${NC}"
    exit 1
fi

# Extract issue details using jq
issue_title=$(echo "$issue_data" | jq -r '.title')
issue_body=$(echo "$issue_data" | jq -r '.body // ""')
issue_state=$(echo "$issue_data" | jq -r '.state')
issue_labels=$(echo "$issue_data" | jq -r '.labels[].name' | paste -sd "," -)
issue_assignees=$(echo "$issue_data" | jq -r '.assignees[].login' | paste -sd "," -)
issue_milestone=$(echo "$issue_data" | jq -r '.milestone.title // ""')

echo -e "${GREEN}✓ Issue details fetched${NC}"
echo -e "  Title: ${issue_title}"
echo -e "  State: ${issue_state}"
echo -e "  Labels: ${issue_labels:-none}"
echo -e "  Assignees: ${issue_assignees:-none}"
echo -e "  Milestone: ${issue_milestone:-none}"
echo ""

# Add reference to source issue in the body
if [ -n "$issue_body" ]; then
    issue_body="${issue_body}

---
_This issue was cloned from ${SOURCE_OWNER}/${SOURCE_REPO}#${ISSUE_NUMBER}_"
else
    issue_body="This issue was cloned from ${SOURCE_OWNER}/${SOURCE_REPO}#${ISSUE_NUMBER}"
fi

# Get all repositories in the target organization
echo "Fetching repositories from ${TARGET_ORG}..."
repo_data=$(gh repo list "$TARGET_ORG" --limit 1000 --json name,isArchived,isFork --jq '.[] | select(.isArchived == false and .isFork == false) | .name')

if [ -z "$repo_data" ]; then
    echo -e "${RED}No repositories found or unable to access organization.${NC}"
    exit 1
fi

repo_count=$(echo "$repo_data" | wc -l)
echo -e "${GREEN}Found ${repo_count} active repositories in ${TARGET_ORG}${NC}"
echo ""

# Ask for confirmation
echo -e "${YELLOW}This will create ${repo_count} new issues across the organization.${NC}"
read -p "Do you want to continue? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi
echo ""

created=0
skipped=0
errors=0

# Process each repository
while IFS= read -r repo; do
    # Skip the source repository to avoid duplicates
    if [ "$TARGET_ORG" = "$SOURCE_OWNER" ] && [ "$repo" = "$SOURCE_REPO" ]; then
        echo -e "${YELLOW}⊘ Skipped: ${TARGET_ORG}/${repo} (source repository)${NC}"
        ((skipped++))
        continue
    fi
    
    echo -e "${BLUE}Processing: ${TARGET_ORG}/${repo}${NC}"
    
    # Check if repository has issues enabled
    repo_info=$(gh repo view "${TARGET_ORG}/${repo}" --json hasIssuesEnabled --jq '.hasIssuesEnabled' 2>&1)
    if [ "$repo_info" != "true" ]; then
        echo -e "${YELLOW}  ⊘ Skipped: Issues not enabled${NC}"
        ((skipped++))
        continue
    fi
    
    # Build gh issue create command arguments
    gh_args=(
        "issue" "create"
        "--repo" "${TARGET_ORG}/${repo}"
        "--title" "$issue_title"
        "--body" "$issue_body"
    )
    
    # Add labels if present
    if [ -n "$issue_labels" ]; then
        gh_args+=("--label" "$issue_labels")
    fi
    
    # Add assignees if present
    if [ -n "$issue_assignees" ]; then
        gh_args+=("--assignee" "$issue_assignees")
    fi
    
    # Add milestone if present
    if [ -n "$issue_milestone" ]; then
        gh_args+=("--milestone" "$issue_milestone")
    fi
    
    # Create the issue
    create_output=$(gh "${gh_args[@]}" 2>&1)
    create_status=$?
    
    if [ $create_status -eq 0 ]; then
        echo -e "${GREEN}  ✓ Issue created: ${create_output}${NC}"
        ((created++))
    else
        echo -e "${RED}  ✗ Failed to create issue${NC}"
        echo -e "${RED}    Error: ${create_output}${NC}"
        ((errors++))
    fi
    
done <<< "$repo_data"

# Summary
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Summary:${NC}"
echo -e "  Total repositories processed: $((created + skipped + errors))"
echo -e "  ${GREEN}Issues created: ${created}${NC}"
echo -e "  ${YELLOW}Skipped: ${skipped}${NC}"
echo -e "  ${RED}Errors: ${errors}${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

if [ $created -gt 0 ]; then
    echo -e "${GREEN}Successfully created ${created} issue(s) across the organization!${NC}"
fi

exit 0
