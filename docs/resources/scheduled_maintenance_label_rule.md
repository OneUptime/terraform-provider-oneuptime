---
page_title: "oneuptime_scheduled_maintenance_label_rule Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure rules for automatically attaching labels to scheduled maintenance events — including labels inherited from the event's monitors — when matching events are created
---

# oneuptime_scheduled_maintenance_label_rule (Resource)

Configure rules for automatically attaching labels to scheduled maintenance events — including labels inherited from the event's monitors — when matching events are created

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this scheduled maintenance label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this scheduled maintenance label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for events on these monitors. Leave empty to match events on any monitor...
- `scheduled_maintenance_labels` (Set) Only trigger for events that already have at least one of these labels. Leave empty to match regardless of event labels...
- `monitor_labels` (Set) Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `title_pattern` (String) Regex (case-insensitive) matched against the event title. Leave empty to match any title...
- `description_pattern` (String) Regex (case-insensitive) matched against the event description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the event when this rule matches. Already-attached labels are not duplicated...
- `inherit_labels_from_monitors` (Bool) When this rule matches, also copy every label of the event's monitors onto the event...
- `inherit_labels_from_hosts` (Bool) When this rule matches, also copy every label of the event's affected hosts onto the event...
- `inherit_labels_from_kubernetes_clusters` (Bool) When this rule matches, also copy every label of the event's affected Kubernetes clusters onto the event...
- `inherit_labels_from_docker_hosts` (Bool) When this rule matches, also copy every label of the event's affected Docker hosts onto the event...
- `inherit_labels_from_podman_hosts` (Bool) When this rule matches, also copy every label of the event's affected Podman hosts onto the event...
- `inherit_labels_from_services` (Bool) When this rule matches, also copy every label of the event's affected services onto the event...
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
terraform import oneuptime_scheduled_maintenance_label_rule.example <id>
```
