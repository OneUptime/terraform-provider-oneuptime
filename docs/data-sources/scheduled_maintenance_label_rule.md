---
page_title: "oneuptime_scheduled_maintenance_label_rule Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure rules for automatically attaching labels to scheduled maintenance events — including labels inherited from the event's monitors — when matching events are created
---

# oneuptime_scheduled_maintenance_label_rule (Data Source)

Configure rules for automatically attaching labels to scheduled maintenance events — including labels inherited from the event's monitors — when matching events are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_label_rule" "by_name" {
  name = "example-scheduled_maintenance_label_rule"
}

data "oneuptime_scheduled_maintenance_label_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this scheduled maintenance label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `monitors` (Set) Only trigger for events on these monitors. Leave empty to match events on any monitor... Computed.
- `scheduled_maintenance_labels` (Set) Only trigger for events that already have at least one of these labels. Leave empty to match regardless of event labels... Computed.
- `monitor_labels` (Set) Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels... Computed.
- `title_pattern` (String) Regex (case-insensitive) matched against the event title. Leave empty to match any title... Computed.
- `description_pattern` (String) Regex (case-insensitive) matched against the event description. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the event when this rule matches. Already-attached labels are not duplicated... Computed.
- `inherit_labels_from_monitors` (Bool) When this rule matches, also copy every label of the event's monitors onto the event... Computed.
- `inherit_labels_from_hosts` (Bool) When this rule matches, also copy every label of the event's affected hosts onto the event... Computed.
- `inherit_labels_from_kubernetes_clusters` (Bool) When this rule matches, also copy every label of the event's affected Kubernetes clusters onto the event... Computed.
- `inherit_labels_from_docker_hosts` (Bool) When this rule matches, also copy every label of the event's affected Docker hosts onto the event... Computed.
- `inherit_labels_from_podman_hosts` (Bool) When this rule matches, also copy every label of the event's affected Podman hosts onto the event... Computed.
- `inherit_labels_from_services` (Bool) When this rule matches, also copy every label of the event's affected services onto the event... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
