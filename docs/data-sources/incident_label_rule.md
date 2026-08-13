---
page_title: "oneuptime_incident_label_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created
---

# oneuptime_incident_label_rule (Data Source)

Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_label_rule" "by_name" {
  name = "example-incident_label_rule"
}

data "oneuptime_incident_label_rule" "by_id" {
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
- `description` (String) Description of this incident label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor... Computed.
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity... Computed.
- `incident_labels` (Set) Only trigger for incidents that already have at least one of these labels. Leave empty to match regardless of incident labels... Computed.
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels... Computed.
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title... Computed.
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the incident when this rule matches. Already-attached labels are not duplicated... Computed.
- `inherit_labels_from_monitors` (Bool) When this rule matches, also copy every label of the incident's monitors onto the incident... Computed.
- `inherit_labels_from_hosts` (Bool) When this rule matches, also copy every label of the incident's affected hosts onto the incident... Computed.
- `inherit_labels_from_kubernetes_clusters` (Bool) When this rule matches, also copy every label of the incident's affected Kubernetes clusters onto the incident... Computed.
- `inherit_labels_from_docker_hosts` (Bool) When this rule matches, also copy every label of the incident's affected Docker hosts onto the incident... Computed.
- `inherit_labels_from_podman_hosts` (Bool) When this rule matches, also copy every label of the incident's affected Podman hosts onto the incident... Computed.
- `inherit_labels_from_services` (Bool) When this rule matches, also copy every label of the incident's affected services onto the incident... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
