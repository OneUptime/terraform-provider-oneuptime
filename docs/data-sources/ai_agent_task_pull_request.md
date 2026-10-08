---
page_title: "oneuptime_ai_agent_task_pull_request Data Source - oneuptime"
subcategory: "Other"
description: |-
  Pull requests created by AI agents during task execution.
---

# oneuptime_ai_agent_task_pull_request (Data Source)

Pull requests created by AI agents during task execution.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai agent task pull request may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_agent_task_pull_request" "example" {
  ai_run_id = "example-ai-run-id"
}

# Or by id:
data "oneuptime_ai_agent_task_pull_request" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_agent_id` (String) ID of the AI Agent that created this pull request. Null when it was proposed from an AI chat conversation. The ID of a `oneuptime_ai_agent`.
- `ai_run_id` (String) ID of the AIRun (runType CodeFix) this pull request was opened by.
- `base_ref_name` (String) The target branch for the pull request.
- `ci_status` (String) Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users.
- `code_repository_id` (String) ID of the Code Repository this pull request was created in. The ID of a `oneuptime_code_repository`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `description` (String) Description/body of the pull request.
- `head_ref_name` (String) The branch name of the pull request (source branch).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `pull_request_id` (Number) The unique ID of the pull request from the hosting platform.
- `pull_request_number` (Number) The pull request number (e.g., #123).
- `pull_request_state` (String) Current state of the pull request (open, closed, merged).
- `pull_request_url` (String) URL to the pull request on the hosting platform.
- `repo_name` (String) Name of the repository.
- `repo_organization_name` (String) Organization or username that owns the repository.
- `runner_verification_status` (String) Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users.
- `runner_verification_summary` (String) Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time.
- `title` (String) Title of the pull request.

### Read-Only

- `ci_status_at` (String) When the CI status last changed.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this pull request belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
