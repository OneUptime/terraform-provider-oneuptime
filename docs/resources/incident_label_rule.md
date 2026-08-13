---
page_title: "oneuptime_incident_label_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created
---

# oneuptime_incident_label_rule (Resource)

Configure rules for automatically attaching labels to incidents — including labels inherited from the incident's monitors and hosts — when matching incidents are created

## Example Usage

```terraform
resource "oneuptime_incident_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this incident label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor...
- `incident_severities` (Set) Only trigger for incidents with these severities. Leave empty to match incidents of any severity...
- `incident_labels` (Set) Only trigger for incidents that already have at least one of these labels. Leave empty to match regardless of incident labels...
- `monitor_labels` (Set) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `incident_title_pattern` (String) Regex (case-insensitive) matched against the incident title. Leave empty to match any title...
- `incident_description_pattern` (String) Regex (case-insensitive) matched against the incident description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the incident when this rule matches. Already-attached labels are not duplicated...
- `inherit_labels_from_monitors` (Bool) When this rule matches, also copy every label of the incident's monitors onto the incident...
- `inherit_labels_from_hosts` (Bool) When this rule matches, also copy every label of the incident's affected hosts onto the incident...
- `inherit_labels_from_kubernetes_clusters` (Bool) When this rule matches, also copy every label of the incident's affected Kubernetes clusters onto the incident...
- `inherit_labels_from_docker_hosts` (Bool) When this rule matches, also copy every label of the incident's affected Docker hosts onto the incident...
- `inherit_labels_from_podman_hosts` (Bool) When this rule matches, also copy every label of the incident's affected Podman hosts onto the incident...
- `inherit_labels_from_services` (Bool) When this rule matches, also copy every label of the incident's affected services onto the incident...
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
terraform import oneuptime_incident_label_rule.example <id>
```
