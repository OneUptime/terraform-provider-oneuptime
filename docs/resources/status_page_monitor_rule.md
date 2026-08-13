---
page_title: "oneuptime_status_page_monitor_rule Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure rules that automatically add matching monitors to a status page group, instead of picking every monitor by hand
---

# oneuptime_status_page_monitor_rule (Resource)

Configure rules that automatically add matching monitors to a status page group, instead of picking every monitor by hand

## Example Usage

```terraform
resource "oneuptime_status_page_monitor_rule" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name of this status page monitor rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this status page monitor rule..
- `is_enabled` (Bool) Whether this rule is enabled. A disabled rule removes the monitors it had added...
- `monitor_labels` (Set) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter...
- `status_page_group_id` (String) A unique identifier for an object, represented as a UUID..
- `show_current_status` (Bool) Show current status like offline, operational or degraded on the resources this rule adds...
- `show_uptime_percent` (Bool) Show uptime percent on the resources this rule adds to the status page...
- `uptime_percent_precision` (String) Precision of the uptime percent shown on the resources this rule adds..
- `show_status_history_chart` (Bool) Show a 90 day uptime history on the resources this rule adds to the status page...
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
terraform import oneuptime_status_page_monitor_rule.example <id>
```
