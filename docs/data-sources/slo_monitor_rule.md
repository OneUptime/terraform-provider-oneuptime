---
page_title: "oneuptime_slo_monitor_rule Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand
---

# oneuptime_slo_monitor_rule (Data Source)

Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one slo monitor rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_slo_monitor_rule" "example" {
  name = "Example slo monitor rule"
}

# Or by id:
data "oneuptime_slo_monitor_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this SLO monitor rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled. A disabled rule matches nothing, so monitors that only this rule attached are detached from the SLO.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor.
- `monitor_type` (String) Only match monitors of this type. Leave empty to skip the type filter.
- `name` (String) Name of this SLO monitor rule.
- `service_level_objective_id` (String) ID of the Service Level Objective this monitor rule belongs to. The ID of a `oneuptime_service_level_objective`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `monitor_labels` (Set of String) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
