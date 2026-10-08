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
  title              = "Example short text"
  pull_request_state = "Example short text"
  description        = "Managed by Terraform"
}
```

## Schema

### Required

- `pull_request_state` (String) Current state of the pull request (open, closed, merged).
- `title` (String) Title of the pull request.

### Optional

- `ai_agent_id` (String) ID of the AI Agent that created this pull request. Null when it was proposed from an AI chat conversation. The ID of a `oneuptime_ai_agent`.
- `ai_run_id` (String) ID of the AIRun (runType CodeFix) this pull request was opened by.
- `base_ref_name` (String) The target branch for the pull request.
- `code_repository_id` (String) ID of the Code Repository this pull request was created in. The ID of a `oneuptime_code_repository`.
- `description` (String) Description/body of the pull request.
- `head_ref_name` (String) The branch name of the pull request (source branch).
- `pull_request_id` (Number) The unique ID of the pull request from the hosting platform.
- `pull_request_number` (Number) The pull request number (e.g., #123).
- `pull_request_url` (String) URL to the pull request on the hosting platform.
- `repo_name` (String) Name of the repository.
- `repo_organization_name` (String) Organization or username that owns the repository.

### Read-Only

- `ci_status` (String) Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users.
- `ci_status_at` (String) When the CI status last changed.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this pull request belongs to. The ID of a `oneuptime_project`.
- `runner_verification_status` (String) Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users.
- `runner_verification_summary` (String) Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ai agent task pull request by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ai_agent_task_pull_request.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_ai_agent_task_pull_request.example <id>
```
