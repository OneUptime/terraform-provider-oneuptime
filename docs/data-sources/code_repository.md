---
page_title: "oneuptime_code_repository Data Source - oneuptime"
subcategory: "Reliability Copilot"
description: |-
  Connect and manage code repositories from GitHub, GitLab, and other providers
---

# oneuptime_code_repository (Data Source)

Connect and manage code repositories from GitHub, GitLab, and other providers Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_code_repository" "by_name" {
  name = "example-code_repository"
}

data "oneuptime_code_repository" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) A description of this code repository.. Computed.
- `repository_hosted_at` (String) Where is this repository hosted (GitHub, GitLab, etc.).. Computed.
- `organization_name` (String) GitHub organization or username that owns this repository.. Computed.
- `repository_name` (String) The name of the repository.. Computed.
- `main_branch_name` (String) The name of the main/default branch.. Computed.
- `setup_command` (String) Command the AI fix Runner executes at the repository root to install dependencies before verifying an AI-authored fix (e.g. 'npm ci'). Runs on your Runner, in the cloned workspace, before the build and test commands. Leave empty to skip... Computed.
- `build_command` (String) Command the AI fix Runner executes at the repository root to verify an AI-authored fix compiles/builds (e.g. 'npm run build'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the build check... Computed.
- `test_command` (String) Command the AI fix Runner executes at the repository root to run the test suite against an AI-authored fix (e.g. 'npm test'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the test check... Computed.
- `max_open_fix_pull_requests` (Number) Maximum AI-authored fix pull requests that may be open on this repository at the same time. At the cap, new AI fix runs are refused a repository token, so they cannot push branches or open pull requests. Unset means the default of 5; 0 blocks AI fix pull requests for this repository entirely... Computed.
- `repository_url` (String) The HTTPS URL to the repository.. Computed.
- `git_hub_app_installation_id` (String) The GitHub App installation ID used to authenticate with this repository.. Computed.
- `git_lab_project_id` (String) The GitLab project ID for this repository.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
