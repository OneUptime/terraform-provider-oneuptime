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
  name        = "Example alert label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert label rule.

### Optional

- `alert_description_pattern` (String) Regex (case-insensitive) matched against the alert description.
- `alert_labels` (Set of String) Only trigger for alerts that already have at least one of these labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only trigger for alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `alert_title_pattern` (String) Regex (case-insensitive) matched against the alert title.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this alert label rule.
- `inherit_labels_from_docker_hosts` (Boolean) When this rule matches, also copy every label of the alert's affected Docker hosts onto the alert. Defaults to `false`.
- `inherit_labels_from_hosts` (Boolean) When this rule matches, also copy every label of the alert's affected hosts onto the alert. Defaults to `false`.
- `inherit_labels_from_kubernetes_clusters` (Boolean) When this rule matches, also copy every label of the alert's affected Kubernetes clusters onto the alert. Defaults to `false`.
- `inherit_labels_from_monitors` (Boolean) When this rule matches, also copy every label of the alert's monitor onto the alert. Defaults to `false`.
- `inherit_labels_from_podman_hosts` (Boolean) When this rule matches, also copy every label of the alert's affected Podman hosts onto the alert. Defaults to `false`.
- `inherit_labels_from_services` (Boolean) When this rule matches, also copy every label of the alert's affected services onto the alert. Defaults to `false`.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the alert when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the alert's monitor description.
- `monitor_labels` (Set of String) Only trigger for alerts from monitors that have at least one of these labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the alert's monitor name.
- `monitors` (Set of String) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_label_rule.example <id>
```
