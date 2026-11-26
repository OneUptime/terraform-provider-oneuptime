---
page_title: "oneuptime_scheduled_maintenance_owner_rule Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure rules for automatically assigning owner users and teams when matching maintenance events are created
---

# oneuptime_scheduled_maintenance_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching maintenance events are created

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_owner_rule" "example" {
  name        = "Example scheduled maintenance owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this scheduled maintenance owner rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this scheduled maintenance owner rule.
- `description_pattern` (String) Regex (case-insensitive) matched against the maintenance event description. Leave empty to match any description.
- `inherit_owners_from_docker_hosts` (Boolean) When this rule matches, also assign every owner of the event's affected Docker hosts to the event. Defaults to `false`.
- `inherit_owners_from_hosts` (Boolean) When this rule matches, also assign every owner of the event's affected hosts to the event. Defaults to `false`.
- `inherit_owners_from_kubernetes_clusters` (Boolean) When this rule matches, also assign every owner of the event's affected Kubernetes clusters to the event. Defaults to `false`.
- `inherit_owners_from_monitors` (Boolean) When this rule matches, also assign every owner of the event's monitors to the event. Defaults to `false`.
- `inherit_owners_from_podman_hosts` (Boolean) When this rule matches, also assign every owner of the event's affected Podman hosts to the event. Defaults to `false`.
- `inherit_owners_from_services` (Boolean) When this rule matches, also assign every owner of the event's affected services to the event. Defaults to `false`.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description.
- `monitor_labels` (Set of String) Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor.
- `monitors` (Set of String) Only trigger for maintenance events on these monitors. Leave empty to match events on any monitor. IDs of `oneuptime_monitor` resources.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the maintenance event when this rule matches. IDs of `oneuptime_team` resources.
- `owner_user_emails` (Set of String) The email addresses of the people in `owner_users`, instead of their IDs: each one of someone in this project, or invited to it. Give one of the two.
- `owner_users` (Set of String) Users to add as owners on the maintenance event when this rule matches. IDs of `oneuptime_user` records. Or name them by email address, in `owner_user_emails`.
- `scheduled_maintenance_labels` (Set of String) Only trigger for events that have at least one of these labels. Leave empty to match regardless of event labels. IDs of `oneuptime_label` resources.
- `title_pattern` (String) Regex (case-insensitive) matched against the maintenance event title. Leave empty to match any title.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_owner_rule.example <id>
```
