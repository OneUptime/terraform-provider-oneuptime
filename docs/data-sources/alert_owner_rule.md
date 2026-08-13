---
page_title: "oneuptime_alert_owner_rule Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically assigning owner users and teams when matching alerts are created
---

# oneuptime_alert_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching alerts are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_owner_rule" "by_name" {
  name = "example-alert_owner_rule"
}

data "oneuptime_alert_owner_rule" "by_id" {
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
- `description` (String) Description of this alert owner rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule.. Computed.
- `monitors` (Set) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor... Computed.
- `alert_severities` (Set) Only trigger for alerts with these severities. Leave empty to match alerts of any severity... Computed.
- `alert_labels` (Set) Only trigger for alerts that have at least one of these labels. Leave empty to match regardless of alert labels... Computed.
- `monitor_labels` (Set) Only trigger for alerts from monitors that have at least one of these labels... Computed.
- `alert_title_pattern` (String) Regex (case-insensitive) matched against the alert title... Computed.
- `alert_description_pattern` (String) Regex (case-insensitive) matched against the alert description... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the alert's monitor name... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the alert's monitor description... Computed.
- `owner_users` (Set) Users to add as owners on the alert when this rule matches... Computed.
- `owner_teams` (Set) Teams to add as owners on the alert when this rule matches... Computed.
- `inherit_owners_from_monitors` (Bool) When this rule matches, also assign every owner of the alert's monitor to the alert... Computed.
- `inherit_owners_from_hosts` (Bool) When this rule matches, also assign every owner of the alert's affected hosts to the alert... Computed.
- `inherit_owners_from_kubernetes_clusters` (Bool) When this rule matches, also assign every owner of the alert's affected Kubernetes clusters to the alert... Computed.
- `inherit_owners_from_docker_hosts` (Bool) When this rule matches, also assign every owner of the alert's affected Docker hosts to the alert... Computed.
- `inherit_owners_from_podman_hosts` (Bool) When this rule matches, also assign every owner of the alert's affected Podman hosts to the alert... Computed.
- `inherit_owners_from_services` (Bool) When this rule matches, also assign every owner of the alert's affected services to the alert... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
