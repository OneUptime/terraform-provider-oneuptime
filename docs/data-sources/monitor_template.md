---
page_title: "oneuptime_monitor_template Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Reusable monitor template. Use it to create new monitors with the same configuration.
---

# oneuptime_monitor_template (Data Source)

Reusable monitor template. Use it to create new monitors with the same configuration. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_monitor_template" "by_name" {
  name = "example-monitor_template"
}

data "oneuptime_monitor_template" "by_id" {
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
- `template_name` (String) Name of the Monitor Template.. Computed.
- `template_description` (String) Description of the Monitor Template.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `monitor_name` (String) Default name applied to monitors created from this template. Users can override on creation. Leave it blank to name each monitor after the resource it watches... Computed.
- `monitor_description` (String) Default description applied to monitors created from this template... Computed.
- `monitor_type` (String) What is the type of monitor created from this template?.. Computed.
- `monitor_steps` (Monitor_steps) MonitorSteps object. Computed.
- `monitoring_interval` (String) Default monitoring interval for monitors created from this template.. Computed.
- `labels` (Set) Default labels applied to monitors created from this template... Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `minimum_probe_agreement` (Number) Default minimum number of probes that must agree on a status before the monitor status changes... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
