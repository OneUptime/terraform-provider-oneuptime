---
page_title: "oneuptime_alert_privacy_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically marking matching alerts as private
---

# oneuptime_alert_privacy_rule (Resource)

Configure rules for automatically marking matching alerts as private

## Example Usage

```terraform
resource "oneuptime_alert_privacy_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this alert privacy rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert privacy rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor...
- `alert_severities` (Set) Only trigger for alerts with these severities. Leave empty to match alerts of any severity...
- `alert_labels` (Set) Only trigger for alerts that have at least one of these labels. Leave empty to match regardless of alert labels...
- `monitor_labels` (Set) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels...
- `alert_title_pattern` (String) Regex (case-insensitive) matched against the alert title. Leave empty to match any title...
- `alert_description_pattern` (String) Regex (case-insensitive) matched against the alert description. Leave empty to match any description...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the alert's monitor name. Leave empty to match any monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the alert's monitor description. Leave empty to match any description...
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
terraform import oneuptime_alert_privacy_rule.example <id>
```
