---
page_title: "oneuptime_ai_agent_task_pull_request Resource - oneuptime"
subcategory: "Other"
description: |-
  Pull requests created by AI agents during task execution.
---

# oneuptime_ai_agent_task_pull_request (Resource)

Pull requests created by AI agents during task execution.

## Example Usage

```terraform
resource "oneuptime_ai_agent_task_pull_request" "example" {
  title = "Example short text"
  pull_request_state = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `title` (String) Title of the pull request...
- `pull_request_state` (String) Current state of the pull request (open, closed, merged)...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `ai_run_id` (String) A unique identifier for an object, represented as a UUID..
- `ai_agent_id` (String) A unique identifier for an object, represented as a UUID..
- `code_repository_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description/body of the pull request...
- `pull_request_url` (String) URL to the pull request on the hosting platform...
- `pull_request_id` (Number) The unique ID of the pull request from the hosting platform...
- `pull_request_number` (Number) The pull request number (e.g., #123)...
- `head_ref_name` (String) The branch name of the pull request (source branch)...
- `base_ref_name` (String) The target branch for the pull request...
- `repo_organization_name` (String) Organization or username that owns the repository...
- `repo_name` (String) Name of the repository...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `ci_status` (String) Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users...
- `ci_status_at` (String) A date time object..
- `runner_verification_status` (String) Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users...
- `runner_verification_summary` (String) Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_ai_agent_task_pull_request.example <id>
```
