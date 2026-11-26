---
page_title: "oneuptime_table_owner_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Rules for automatically assigning owners to tables when matching resources are created.
---

# oneuptime_table_owner_rule (Data Source)

Rules for automatically assigning owners to tables when matching resources are created.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one table owner rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_table_owner_rule" "example" {
  name = "Example table owner rule"
}

# Or by id:
data "oneuptime_table_owner_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `database_regex_pattern` (String) Case-insensitive regex matched against the name of the database the table is in.
- `description` (String) Description of this rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of this rule.
- `name_regex_pattern` (String) Case-insensitive regex matched against the table's name.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule.
- `schema_regex_pattern` (String) Case-insensitive regex matched against the table's schema.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `match_labels` (Set of String) Only tables with one of these labels. Leave empty to skip the filter. IDs of `oneuptime_label` resources.
- `owner_teams` (Set of String) Teams to add as owners on the resource when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the resource when this rule matches. IDs of `oneuptime_user` records.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
