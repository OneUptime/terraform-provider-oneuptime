---
page_title: "oneuptime_team_or_label_budget Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Daily ingestion budgets for the telemetry of the services a team owns, or that carry a label.
---

# oneuptime_team_or_label_budget (Data Source)

Daily ingestion budgets for the telemetry of the services a team owns, or that carry a label.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team or label budget may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_or_label_budget" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_or_label_budget" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `label_id` (String) ID of the label whose services the budget is for. Empty for a team's budget. The ID of a `oneuptime_label`.
- `logs_daily_ingestion_budget_in_gb` (Number) The most logs the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their logs are dropped until the budget resets. At most the project's own budget. Empty: no logs budget of their own (the project's still applies).
- `metrics_daily_ingestion_budget_in_gb` (Number) The most metrics the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their metrics are dropped until the budget resets. At most the project's own budget. Empty: no metrics budget of their own (the project's still applies).
- `profiles_daily_ingestion_budget_in_gb` (Number) The most profiles the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their profiles are dropped until the budget resets. At most the project's own budget. Empty: no profiles budget of their own (the project's still applies).
- `team_id` (String) ID of the team whose services the budget is for. Empty for a label's budget. The ID of a `oneuptime_team`.
- `traces_daily_ingestion_budget_in_gb` (Number) The most traces the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their traces are dropped until the budget resets. At most the project's own budget. Empty: no traces budget of their own (the project's still applies).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
