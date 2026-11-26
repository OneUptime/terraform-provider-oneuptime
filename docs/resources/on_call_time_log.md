---
page_title: "oneuptime_on_call_time_log Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  When each person was on call for an on-call policy, from when their shift started to when it ended: what the User On-Call Time report adds up.
---

# oneuptime_on_call_time_log (Resource)

When each person was on call for an on-call policy, from when their shift started to when it ended: what the User On-Call Time report adds up.

## Example Usage

```terraform
resource "oneuptime_on_call_time_log" "example" {
  starts_at  = "2030-01-01T00:00:00Z"
  user_email = "alice@example.com"
}
```

## Schema

### Required

- `starts_at` (String) When does this start?

### Optional

- `ends_at` (String) When does this end? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `more_info` (String) More information about this log record.
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule ID where this escalation rule belongs.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Schedule where this escalation rule belongs.
- `team_id` (String) ID of the team this log belongs to.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User ID for which this log belongs. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

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

Import an existing on call time log by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_time_log.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_time_log.example <id>
```
