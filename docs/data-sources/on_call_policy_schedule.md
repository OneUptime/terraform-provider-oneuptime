---
page_title: "oneuptime_on_call_policy_schedule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Who is on call at any moment: people take turns in layers, and each one hands over to the next.
---

# oneuptime_on_call_policy_schedule (Data Source)

Who is on call at any moment: people take turns in layers, and each one hands over to the next.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call policy schedule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_policy_schedule" "example" {
  name = "Example on call policy schedule"
}

# Or by id:
data "oneuptime_on_call_policy_schedule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_user_id_on_roster` (String) User ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `next_user_id_on_roster` (String) Next ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).
- `shift_config_version` (Number) Incremented whenever the schedule's layers, members, overrides or policy attachments change. Used as the calendar feed SEQUENCE.
- `slug` (String) Friendly globally unique name for your object.
- `timezone` (String) IANA timezone this schedule's restriction and hand-off wall-clock times are interpreted in. When empty, times are interpreted in the server's local timezone (legacy behavior).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `roster_handoff_at` (String) When does the roster handoff occur for this schedule for the current user?
- `roster_next_handoff_at` (String) When does the next roster handoff occur for this schedule for the next user?
- `roster_next_start_at` (String) When does the next event start for this schedule for the next user?
- `roster_start_at` (String) When does the current event start for this schedule for the current user?
- `updated_at` (String) Date and Time when the object was updated.
