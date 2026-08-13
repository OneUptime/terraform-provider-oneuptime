---
page_title: "oneuptime_code_repository Resource - oneuptime"
subcategory: "Reliability Copilot"
description: |-
  Connect and manage code repositories from GitHub, GitLab, and other providers
---

# oneuptime_code_repository (Resource)

Connect and manage code repositories from GitHub, GitLab, and other providers

## Example Usage

```terraform
resource "oneuptime_code_repository" "example" {
  name = "Example short text"
  repository_hosted_at = "Example short text"
  organization_name = "Example short text"
  repository_name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) A friendly name for this code repository..
- `repository_hosted_at` (String) Where is this repository hosted (GitHub, GitLab, etc.)..
- `organization_name` (String) GitHub organization or username that owns this repository..
- `repository_name` (String) The name of the repository..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) A description of this code repository..
- `main_branch_name` (String) The name of the main/default branch..
- `setup_command` (String) Command the AI fix Runner executes at the repository root to install dependencies before verifying an AI-authored fix (e.g. 'npm ci'). Runs on your Runner, in the cloned workspace, before the build and test commands. Leave empty to skip...
- `build_command` (String) Command the AI fix Runner executes at the repository root to verify an AI-authored fix compiles/builds (e.g. 'npm run build'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the build check...
- `test_command` (String) Command the AI fix Runner executes at the repository root to run the test suite against an AI-authored fix (e.g. 'npm test'). A failure is fed back to the code agent for bounded repair attempts before the pull request opens. Leave empty to skip the test check...
- `max_open_fix_pull_requests` (Number) Maximum AI-authored fix pull requests that may be open on this repository at the same time. At the cap, new AI fix runs are refused a repository token, so they cannot push branches or open pull requests. Unset means the default of 5; 0 blocks AI fix pull requests for this repository entirely...
- `repository_url` (String) The HTTPS URL to the repository..
- `git_lab_project_id` (String) The GitLab project ID for this repository..
- `secret_token` (String) Secret token used to verify incoming webhooks..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `git_hub_app_installation_id` (String) The GitHub App installation ID used to authenticate with this repository..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_code_repository.example <id>
```
