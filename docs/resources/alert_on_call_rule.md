---
page_title: "oneuptime_alert_on_call_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically executing on-call duty policies when matching alerts are created
---

# oneuptime_alert_on_call_rule (Resource)

Configure rules for automatically executing on-call duty policies when matching alerts are created

## Example Usage

```terraform
resource "oneuptime_alert_on_call_rule" "example" {
  name        = "Example alert on call rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert on-call rule.

### Optional

- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description.
- `alert_labels` (Set of String) Only trigger for alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only trigger for alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this alert on-call rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_labels` (Set of String) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `monitors` (Set of String) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) On-call duty policies to execute when an alert matches this rule. IDs of `oneuptime_on_call_policy` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert on call rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_on_call_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_on_call_rule.example <id>
```
