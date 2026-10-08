---
page_title: "oneuptime_scheduled_maintenance_owner_rule Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure rules for automatically assigning owner users and teams when matching scheduled maintenance events are created
---

# oneuptime_scheduled_maintenance_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching scheduled maintenance events are created

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
- `description_pattern` (String) Regex (case-insensitive) matched against the scheduled maintenance event description. Leave empty to match any description.
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
- `monitors` (Set of String) Only trigger for scheduled maintenance events on these monitors. Leave empty to match events on any monitor. IDs of `oneuptime_monitor` resources.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the scheduled maintenance event when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the scheduled maintenance event when this rule matches. IDs of `oneuptime_user` records.
- `scheduled_maintenance_labels` (Set of String) Only trigger for events that have at least one of these labels. Leave empty to match regardless of event labels. IDs of `oneuptime_label` resources.
- `title_pattern` (String) Regex (case-insensitive) matched against the scheduled maintenance event title. Leave empty to match any title.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
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
