---
page_title: "oneuptime_alert_label_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically attaching labels to alerts — including labels inherited from the alert's monitor — when matching alerts are created
---

# oneuptime_alert_label_rule (Resource)

Configure rules for automatically attaching labels to alerts — including labels inherited from the alert's monitor — when matching alerts are created

## Example Usage

```terraform
resource "oneuptime_alert_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this alert label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor...
- `alert_severities` (Set) Only trigger for alerts with these severities. Leave empty to match alerts of any severity...
- `alert_labels` (Set) Only trigger for alerts that already have at least one of these labels...
- `monitor_labels` (Set) Only trigger for alerts from monitors that have at least one of these labels...
- `alert_title_pattern` (String) Regex (case-insensitive) matched against the alert title...
- `alert_description_pattern` (String) Regex (case-insensitive) matched against the alert description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the alert's monitor name...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the alert's monitor description...
- `labels_to_add` (Set) Labels to attach to the alert when this rule matches. Already-attached labels are not duplicated...
- `inherit_labels_from_monitors` (Bool) When this rule matches, also copy every label of the alert's monitor onto the alert...
- `inherit_labels_from_hosts` (Bool) When this rule matches, also copy every label of the alert's affected hosts onto the alert...
- `inherit_labels_from_kubernetes_clusters` (Bool) When this rule matches, also copy every label of the alert's affected Kubernetes clusters onto the alert...
- `inherit_labels_from_docker_hosts` (Bool) When this rule matches, also copy every label of the alert's affected Docker hosts onto the alert...
- `inherit_labels_from_podman_hosts` (Bool) When this rule matches, also copy every label of the alert's affected Podman hosts onto the alert...
- `inherit_labels_from_services` (Bool) When this rule matches, also copy every label of the alert's affected services onto the alert...
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
terraform import oneuptime_alert_label_rule.example <id>
```
