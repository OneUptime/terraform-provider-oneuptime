---
page_title: "oneuptime_incident_label_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created
---

# oneuptime_incident_label_rule (Data Source)

Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_label_rule" "example" {
  name = "Example incident label rule"
}

# Or by id:
data "oneuptime_incident_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this incident label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description.
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title.
- `inherit_labels_from_docker_hosts` (Boolean) When this rule matches, also copy every label of the incident's affected Docker hosts onto the incident.
- `inherit_labels_from_hosts` (Boolean) When this rule matches, also copy every label of the incident's affected hosts onto the incident.
- `inherit_labels_from_kubernetes_clusters` (Boolean) When this rule matches, also copy every label of the incident's affected Kubernetes clusters onto the incident.
- `inherit_labels_from_monitors` (Boolean) When this rule matches, also copy every label of the incident's monitors onto the incident.
- `inherit_labels_from_podman_hosts` (Boolean) When this rule matches, also copy every label of the incident's affected Podman hosts onto the incident.
- `inherit_labels_from_services` (Boolean) When this rule matches, also copy every label of the incident's affected services onto the incident.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor.
- `name` (String) Name of this incident label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `incident_labels` (Set of String) Only trigger for incidents that already have at least one of these labels. Leave empty to match regardless of incident labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only trigger for incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `labels_to_add` (Set of String) Labels to attach to the incident when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
