---
page_title: "oneuptime_database_label_rule Data Source - oneuptime"
subcategory: "Infrastructure"
description: |-
  Configure rules for automatically attaching labels to databases when matching databases are created
---

# oneuptime_database_label_rule (Data Source)

Configure rules for automatically attaching labels to databases when matching databases are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one database label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_database_label_rule" "example" {
  name = "Example database label rule"
}

# Or by id:
data "oneuptime_database_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `database_server_description_pattern` (String) Regex (case-insensitive) matched against the database description. Leave empty to match any description.
- `database_server_name_pattern` (String) Regex (case-insensitive) matched against the database name. Discovered databases are named after their engine (e.g. PostgreSQL orders-db:5432), so ^PostgreSQL matches every PostgreSQL database. Leave empty to match any name.
- `description` (String) Description of this database label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of this database label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `database_server_labels` (Set of String) Only trigger for databases that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `labels_to_add` (Set of String) Labels to attach to the database when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
