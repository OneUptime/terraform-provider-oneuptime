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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this scheduled maintenance owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this scheduled maintenance owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `monitors` (Set) Only trigger for scheduled maintenance events on these monitors. Leave empty to match events on any monitor...
- `scheduled_maintenance_labels` (Set) Only trigger for events that have at least one of these labels. Leave empty to match regardless of event labels...
- `monitor_labels` (Set) Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `title_pattern` (String) Regex (case-insensitive) matched against the scheduled maintenance event title. Leave empty to match any title...
- `description_pattern` (String) Regex (case-insensitive) matched against the scheduled maintenance event description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the scheduled maintenance event when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the scheduled maintenance event when this rule matches...
- `inherit_owners_from_monitors` (Bool) When this rule matches, also assign every owner of the event's monitors to the event...
- `inherit_owners_from_hosts` (Bool) When this rule matches, also assign every owner of the event's affected hosts to the event...
- `inherit_owners_from_kubernetes_clusters` (Bool) When this rule matches, also assign every owner of the event's affected Kubernetes clusters to the event...
- `inherit_owners_from_docker_hosts` (Bool) When this rule matches, also assign every owner of the event's affected Docker hosts to the event...
- `inherit_owners_from_podman_hosts` (Bool) When this rule matches, also assign every owner of the event's affected Podman hosts to the event...
- `inherit_owners_from_services` (Bool) When this rule matches, also assign every owner of the event's affected services to the event...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_owner_rule.example <id>
```
