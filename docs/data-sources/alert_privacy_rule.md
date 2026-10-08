---
page_title: "oneuptime_alert_privacy_rule Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically marking matching alerts as private
---

# oneuptime_alert_privacy_rule (Data Source)

Configure rules for automatically marking matching alerts as private

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert privacy rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_privacy_rule" "example" {
  name = "Example alert privacy rule"
}

# Or by id:
data "oneuptime_alert_privacy_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_description_pattern` (String) Regex (case-insensitive) matched against the alert description. Leave empty to match any description.
- `alert_title_pattern` (String) Regex (case-insensitive) matched against the alert title. Leave empty to match any title.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this alert privacy rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the alert's monitor description. Leave empty to match any description.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the alert's monitor name. Leave empty to match any monitor.
- `name` (String) Name of this alert privacy rule.

### Read-Only

- `alert_labels` (Set of String) Only trigger for alerts that have at least one of these labels. Leave empty to match regardless of alert labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only trigger for alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `monitor_labels` (Set of String) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
