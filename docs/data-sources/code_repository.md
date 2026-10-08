---
page_title: "oneuptime_code_repository Data Source - oneuptime"
subcategory: "Reliability Copilot"
description: |-
  Connect and manage code repositories from GitHub, GitLab, and other providers
---

# oneuptime_code_repository (Data Source)

Connect and manage code repositories from GitHub, GitLab, and other providers

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one code repository may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_code_repository" "example" {
  name = "Example code repository"
}

# Or by id:
data "oneuptime_code_repository" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `build_command` (String) Command the AI fix Runner executes at the repository root to verify an AI-authored fix compiles/builds (e.g. 'npm run build'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the build check.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) A description of this code repository.
- `git_hub_app_installation_id` (String) The GitHub App installation ID used to authenticate with this repository.
- `git_hub_trigger_label` (String) The issue label that hands an issue to the OneUptime GitHub App. Adding this label to an issue starts the same work an '@mention implement this' would, which is how you assign work to the app from the GitHub UI. Unset means 'oneuptime'.
- `git_lab_project_id` (String) The GitLab project ID for this repository.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_git_hub_commands_enabled` (Boolean) Whether the OneUptime GitHub App acts on mentions, assignments and trigger labels in this repository. Only people with write access to the repository can command it, and it never merges anything. Unset means enabled.
- `main_branch_name` (String) The name of the main/default branch.
- `max_open_fix_pull_requests` (Number) Maximum AI-authored fix pull requests that may be open on this repository at the same time. At the cap, new AI fix runs are refused a repository token, so they cannot push branches or open pull requests. Unset means no cap; 0 blocks AI fix pull requests for this repository entirely.
- `name` (String) A friendly name for this code repository.
- `organization_name` (String) GitHub organization or username that owns this repository.
- `repository_hosted_at` (String) Where is this repository hosted (GitHub, GitLab, etc.).
- `repository_name` (String) The name of the repository.
- `repository_url` (String) The HTTPS URL to the repository.
- `setup_command` (String) Command the AI fix Runner executes at the repository root to install dependencies before verifying an AI-authored fix (e.g. 'npm ci'). Runs on your Runner, in the cloned workspace, before the build and test commands. Leave empty to skip.
- `slug` (String) Friendly globally unique name for your object.
- `test_command` (String) Command the AI fix Runner executes at the repository root to run the test suite against an AI-authored fix (e.g. 'npm test'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the test check.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
