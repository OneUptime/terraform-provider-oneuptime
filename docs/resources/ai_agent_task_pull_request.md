---
page_title: "oneuptime_ai_agent_task_pull_request Resource - oneuptime"
subcategory: "AI"
description: |-
  Pull requests created by AI agents during task execution.
---

# oneuptime_ai_agent_task_pull_request (Resource)

Pull requests created by AI agents during task execution.

## Example Usage

```terraform
resource "oneuptime_ai_agent_task_pull_request" "example" {
  title              = "Example short text"
  pull_request_state = "open"
  description        = "Managed by Terraform"
}
```

## Schema

### Required

- `pull_request_state` (String) Current state of the pull request (open, closed, merged). OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Allowed values: `open`, `closed`, `merged`.
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

- `carried_alert_ids` (String) IDs of the alerts this pull request carries: the subject of the AI chat it was opened from. A JSON value: read it with `jsondecode()`.
- `carried_incident_ids` (String) IDs of the incidents this pull request carries: the subject of the AI chat it was opened from. A JSON value: read it with `jsondecode()`.
- `ci_status` (String) Rolled-up conclusion of the repository's own CI check runs on this pull request (Pending, Green, Red, ExpectedFailureObserved for should-fail regression-test PRs, NoCiConfigured). Null until the sync job first polls check runs. Written by AIAgent:SyncPullRequestStates — never by users.
- `ci_status_at` (String) When the CI status last changed.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `deployed_at` (String) When the first recorded deploy that contains the fix ran. Written by the server.
- `deployed_commit_sha` (String) The commit that deploy ran: the fix's merge commit or a later one. Written by the server.
- `deployed_version` (String) The version that deploy reported, when it reported one. Written by the server.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `merge_commit_sha` (String) The commit GitHub merged the pull request as. The fix is in production once a recorded deploy contains this commit. Written by the server when the pull request merges.
- `merged_at` (String) When GitHub merged the pull request. Written by the server.
- `production_check` (String) What production is watched for once the fix merges: Exception (the exception stops), Performance (the operation it names gets faster), Recurrence (the monitor stops raising incidents or alerts) or Telemetry (the spans and log fields it adds arrive). Empty on a watched fix of an exception from before the column existed. Written by the server when the pull request merges.
- `production_status` (String) Whether the merged fix worked in production: AwaitingDeploy, Watching, Verified, Regressed or Unverifiable. Written by the server, which watches what the fix was for (Production Check) once a recorded deploy contains the fix.
- `production_status_at` (String) When the production status last changed. Written by the server.
- `production_summary` (String) What the production check found, in a sentence: the deploy that shipped the fix and what it measured before and since. Written by the server.
- `project_id` (String) ID of the project this pull request belongs to. The ID of a `oneuptime_project`.
- `runner_verification_status` (String) Outcome of the Runner-side build/test verification that ran against the fix BEFORE this pull request opened (Passed, Failed, Skipped when the repository has no verification commands configured). Distinct from CI Status, which mirrors the repository's own CI checks after the PR exists. Written by the Runner at record time — never by users.
- `runner_verification_summary` (String) Human-readable summary of the Runner-side verification (which commands ran, what failed, how many repair attempts were used). Written by the Runner at record time.
- `security_case_id` (String) ID of the security case this pull request fixes, when Fix With AI opened it from the case.
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
