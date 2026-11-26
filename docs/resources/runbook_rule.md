---
page_title: "oneuptime_runbook_rule Resource - oneuptime"
subcategory: "Runbooks & Runners"
description: |-
  Auto-attach runbooks to incidents, alerts, or maintenance events when they are created.
---

# oneuptime_runbook_rule (Resource)

Auto-attach runbooks to incidents, alerts, or maintenance events when they are created.

## Example Usage

```terraform
resource "oneuptime_runbook_rule" "example" {
  name                = "Example runbook rule"
  trigger_entity_type = "Incident"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this runbook rule.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance. Allowed values: `Incident`, `Alert`, `ScheduledMaintenance`.

### Optional

- `alert_severities` (Set of String) Only match alerts with one of these severities. Alert rules only. Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this runbook rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `incident_severities` (Set of String) Only match incidents with one of these severities. Incident rules only. Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels` (Set of String) Only match incidents, alerts or maintenance events that carry at least one of these labels. Leave empty to match regardless of their labels. IDs of `oneuptime_label` resources.
- `monitor_description_pattern` (String) Case-insensitive regex matched against the descriptions of the monitors of the incident, alert or maintenance event. Leave empty to match any monitor description.
- `monitor_labels` (Set of String) Only match when a monitor of the incident, alert or maintenance event carries at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Case-insensitive regex matched against the names of the monitors of the incident, alert or maintenance event. Leave empty to match any monitor name.
- `monitors` (Set of String) Only match incidents and maintenance events that affect, and alerts raised by, at least one of these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `runbooks` (Set of String) Runbooks to start when this rule matches. Each runbook produces its own execution. IDs of `oneuptime_runbook` resources.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook_rule.example <id>
```
