---
page_title: "oneuptime_data_security_issue Data Source - oneuptime"
subcategory: "Other"
description: |-
  Everything wrong with a store that a check found, ranked by risk: what is wrong, why it matters, how to fix it, and who owns it.
---

# oneuptime_data_security_issue (Data Source)

Everything wrong with a store that a check found, ranked by risk: what is wrong, why it matters, how to fix it, and who owns it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security issue may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_issue" "example" {
  data_store_id = oneuptime_data_store.example.id
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
- `alert_id` (String) The alert raised for the issue, if any. The ID of a `oneuptime_alert`.
- `assigned_to_team_id` (String) The team working on the issue. The ID of a `oneuptime_team`.
- `assigned_to_user_id` (String) The person working on the issue. The ID of a `oneuptime_user` (see the data source).
- `check_group` (String) The check's group: Exposure, Encryption, Recovery, Visibility, Access, Data hygiene, Telemetry, Residency or AI.
- `check_id` (String) The check that failed.
- `data_store_id` (String) The store the issue is on. The ID of a `oneuptime_data_store`.
- `evidence` (String) What the check saw, in one sentence.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) The incident declared for the issue, if any. The ID of a `oneuptime_incident`.
- `is_owner_notified` (Boolean) Whether the store's owners were told about the issue.
- `issue_key` (String) The store and check, which identify the issue within the project.
- `reopen_count` (Number) How many times the issue opened again after it was resolved.
- `risk_score` (Number) Orders issues inside a severity, from 0 to 100: severity, exposure, volume, recent activity and production.
- `security_case_id` (String) The security case opened for the issue, if any. The ID of a `oneuptime_security_case`.
- `severity` (String) How severe the issue is for this store's data: Critical, High, Medium or Low.
- `status` (String) Open, Accepted, Resolved or Not an issue.
- `title` (String) What is wrong, in a few words.

### Read-Only

- `accepted_until` (String) When an accepted issue opens again if its check still fails.
- `created_at` (String) Date and Time when the object was created.
- `data_type_ids` (String) The sensitive data types the store holds, which make the issue as severe as it is. A JSON value: write it with `jsonencode()`.
- `first_detected_at` (String) When the check first failed.
- `frameworks` (String) The framework requirements the check is evidence for. A JSON value: write it with `jsonencode()`.
- `last_detected_at` (String) When the check last failed.
- `project_id` (String) ID of the project this issue belongs to. The ID of a `oneuptime_project`.
- `recheck_requested_at` (String) When a person last asked for the store to be checked again.
- `resolved_at` (String) When the check last passed after failing.
- `risk_factors` (String) Each part of the risk score, with its points. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
