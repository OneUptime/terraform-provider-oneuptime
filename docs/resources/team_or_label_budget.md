---
page_title: "oneuptime_team_or_label_budget Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Daily ingestion budgets for the telemetry of the services, hosts and other resources a team owns, or that carry a label.
---

# oneuptime_team_or_label_budget (Resource)

Daily ingestion budgets for the telemetry of the services, hosts and other resources a team owns, or that carry a label.

## Example Usage

```terraform
resource "oneuptime_team_or_label_budget" "example" {

}
```

## Schema

### Optional

- `label_id` (String) ID of the label the budget is for: every service, host and other resource that carries it. Empty for a team's budget. The ID of a `oneuptime_label`.
- `logs_daily_ingestion_budget_in_gb` (Number) The most logs the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their logs are dropped until the budget resets. At most the project's own budget. Empty: no logs budget of their own (the project's still applies).
- `metrics_daily_ingestion_budget_in_gb` (Number) The most metrics the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their metrics are dropped until the budget resets. At most the project's own budget. Empty: no metrics budget of their own (the project's still applies).
- `profiles_daily_ingestion_budget_in_gb` (Number) The most profiles the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their profiles are dropped until the budget resets. At most the project's own budget. Empty: no profiles budget of their own (the project's still applies).
- `team_id` (String) ID of the team the budget is for: every service, host and other resource it owns. Empty for a label's budget. The ID of a `oneuptime_team`.
- `traces_daily_ingestion_budget_in_gb` (Number) The most traces the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their traces are dropped until the budget resets. At most the project's own budget. Empty: no traces budget of their own (the project's still applies).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team or label budget by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_or_label_budget.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_or_label_budget.example <id>
```
