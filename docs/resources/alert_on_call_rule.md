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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this alert on-call rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert on-call rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `monitors` (Set) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor...
- `alert_severities` (Set) Only trigger for alerts with these severities. Leave empty to match alerts of any severity...
- `alert_labels` (Set) Only trigger for alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels...
- `monitor_labels` (Set) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels...
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'...
- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description...
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'...
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description...
- `on_call_duty_policies` (Set) On-call duty policies to execute when an alert matches this rule...
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
terraform import oneuptime_alert_on_call_rule.example <id>
```
