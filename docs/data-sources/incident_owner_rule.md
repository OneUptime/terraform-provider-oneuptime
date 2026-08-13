---
page_title: "oneuptime_incident_owner_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically assigning owner users and teams when matching incidents are created
---

# oneuptime_incident_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching incidents are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_owner_rule" "by_name" {
  name = "example-incident_owner_rule"
}

data "oneuptime_incident_owner_rule" "by_id" {
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
- `description` (String) Description of this incident owner rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule.. Computed.
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor... Computed.
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity... Computed.
- `incident_labels` (Set) Only trigger for incidents that have at least one of these labels. Leave empty to match regardless of incident labels... Computed.
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels... Computed.
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title... Computed.
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description... Computed.
- `owner_users` (Set) Users to add as owners on the incident when this rule matches... Computed.
- `owner_teams` (Set) Teams to add as owners on the incident when this rule matches... Computed.
- `inherit_owners_from_monitors` (Bool) When this rule matches, also assign every owner of the incident's monitors to the incident... Computed.
- `inherit_owners_from_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected hosts to the incident... Computed.
- `inherit_owners_from_kubernetes_clusters` (Bool) When this rule matches, also assign every owner of the incident's affected Kubernetes clusters to the incident... Computed.
- `inherit_owners_from_docker_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected Docker hosts to the incident... Computed.
- `inherit_owners_from_podman_hosts` (Bool) When this rule matches, also assign every owner of the incident's affected Podman hosts to the incident... Computed.
- `inherit_owners_from_services` (Bool) When this rule matches, also assign every owner of the incident's affected services to the incident... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
