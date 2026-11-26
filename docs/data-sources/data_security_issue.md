---
page_title: "oneuptime_data_security_issue Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Everything wrong with a store that a check found, ranked by risk: what is wrong, why it matters, how to fix it, and who owns it.
---

# oneuptime_data_security_issue (Data Source)

Everything wrong with a store that a check found, ranked by risk: what is wrong, why it matters, how to fix it, and who owns it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security issue may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_issue" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  assigned_to_user_email   = "alice@example.com"
}

# Or by id:
data "oneuptime_data_security_issue" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `accepted_by_user_id` (String) The person who accepted the risk. The ID of a `oneuptime_user` (see the data source).
- `accepted_reason` (String) Why the risk was accepted.
- `ai_analysis_llm_log_id` (String) The AI Logs row of the model call that wrote the analysis: which model, its tokens and its cost.
- `ai_verdict` (String) OneUptime AI's verdict on the issue: Real, Likely by design, Likely a wrong finding, or Cannot tell. A person decides what to do with it.
- `alert_id` (String) The alert raised for the issue, if any. The ID of a `oneuptime_alert`.
- `assigned_to_team_id` (String) The team working on the issue. The ID of a `oneuptime_team`.
- `assigned_to_user_email` (String) Look it up by the email address of the person in `assigned_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `assigned_to_user_id` (String) The person working on the issue. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `assigned_to_user_email`.
- `check_group` (String) The check's group: Exposure, Encryption, Recovery, Visibility, Access, Data hygiene, Telemetry, Residency or AI.
- `check_id` (String) The check that failed.
- `code_fix_task_id` (String) The Code task whose pull request fixes this issue in code: Fix with AI's, or the one that stops its service writing its telemetry leaks. Empty until a fix is asked for.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_security_issue_group_id` (String) The group of issues of the same check whose stores share what made them, so the cause is fixed once: one pull request, or one approval. Empty for an issue in no group.
- `data_store_id` (String) The store the issue is on. The ID of a `oneuptime_data_store`.
- `evidence` (String) What the check saw, in one sentence.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) The incident declared for the issue, if any. The ID of a `oneuptime_incident`.
- `is_owner_notified` (Boolean) Whether the store's owners were told about the issue.
- `issue_key` (String) The store and check, which identify the issue within the project.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `not_an_issue_reason` (String) Why a person said the finding was wrong.
- `reopen_count` (Number) How many times the issue opened again after it was resolved.
- `resolution` (String) How a resolved issue was resolved: FixedInCode, FixedInInfrastructureCode, FixedFromHere, ContainedAtIngest (a scrub rule hides the data at ingest, but the code still writes it), CheckPassed, ResolvedByPerson, NoLongerChecked or StoreArchived. Empty while the issue is not resolved.
- `risk_score` (Number) Orders issues inside a severity, from 0 to 100: severity, exposure, volume, recent activity and production.
- `security_case_id` (String) The security case opened for the issue, if any. The ID of a `oneuptime_security_case`.
- `severity` (String) How severe the issue is for this store's data: Critical, High, Medium or Low.
- `status` (String) Open, Accepted, Resolved or Not an issue.
- `title` (String) What is wrong, in a few words.

### Read-Only

- `accepted_until` (String) When an accepted issue opens again if its check still fails.
- `ai_analyzed_at` (String) When OneUptime AI last wrote its analysis of the issue.
- `created_at` (String) Date and Time when the object was created.
- `data_type_ids` (String) The sensitive data types the store holds, which make the issue as severe as it is. A JSON value: read it with `jsondecode()`.
- `first_detected_at` (String) When the check first failed.
- `frameworks` (String) The framework requirements the check is evidence for. A JSON value: read it with `jsondecode()`.
- `last_detected_at` (String) When the check last failed.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this issue belongs to. The ID of a `oneuptime_project`.
- `recheck_requested_at` (String) When a person last asked for the store to be checked again.
- `resolution_details` (String) The facts behind the resolution: the pull request that fixed it, the version that shipped it and when production verified it; the scrub rules that contain it; the response action that fixed it; or who resolved it. A JSON value: read it with `jsondecode()`.
- `resolved_at` (String) When the check last passed after failing.
- `risk_factors` (String) Each part of the risk score, with its points. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
