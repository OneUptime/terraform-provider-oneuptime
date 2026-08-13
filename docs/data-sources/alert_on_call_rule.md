---
page_title: "oneuptime_alert_on_call_rule Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically executing on-call duty policies when matching alerts are created
---

# oneuptime_alert_on_call_rule (Data Source)

Configure rules for automatically executing on-call duty policies when matching alerts are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_on_call_rule" "by_name" {
  name = "example-alert_on_call_rule"
}

data "oneuptime_alert_on_call_rule" "by_id" {
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
- `description` (String) Description of this alert on-call rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `monitors` (Set) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor... Computed.
- `alert_severities` (Set) Only trigger for alerts with these severities. Leave empty to match alerts of any severity... Computed.
- `alert_labels` (Set) Only trigger for alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels... Computed.
- `monitor_labels` (Set) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels... Computed.
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'... Computed.
- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'... Computed.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description... Computed.
- `on_call_duty_policies` (Set) On-call duty policies to execute when an alert matches this rule... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
