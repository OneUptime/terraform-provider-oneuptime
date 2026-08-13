---
page_title: "oneuptime_ai_agent_task_pull_request Data Source - oneuptime"
subcategory: "Other"
description: |-
  Pull requests created by AI agents during task execution.
---

# oneuptime_ai_agent_task_pull_request (Data Source)

Pull requests created by AI agents during task execution. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_agent_task_pull_request" "by_name" {
  name = "example-ai_agent_task_pull_request"
}

data "oneuptime_ai_agent_task_pull_request" "by_id" {
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
- `ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `ai_agent_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `code_repository_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `title` (String) Title of the pull request... Computed.
- `description` (String) Description/body of the pull request... Computed.
- `pull_request_url` (String) URL to the pull request on the hosting platform... Computed.
- `pull_request_id` (Number) The unique ID of the pull request from the hosting platform... Computed.
- `pull_request_number` (Number) The pull request number (e.g., #123)... Computed.
- `pull_request_state` (String) Current state of the pull request (open, closed, merged)... Computed.
- `ci_status` (String) Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users... Computed.
- `ci_status_at` (String) A date time object.. Computed.
- `runner_verification_status` (String) Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users... Computed.
- `runner_verification_summary` (String) Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time... Computed.
- `head_ref_name` (String) The branch name of the pull request (source branch)... Computed.
- `base_ref_name` (String) The target branch for the pull request... Computed.
- `repo_organization_name` (String) Organization or username that owns the repository... Computed.
- `repo_name` (String) Name of the repository... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
