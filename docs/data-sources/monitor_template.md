---
page_title: "oneuptime_monitor_template Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Reusable monitor template. Use it to create new monitors with the same configuration.
---

# oneuptime_monitor_template (Data Source)

Reusable monitor template. Use it to create new monitors with the same configuration.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_template" "example" {
  template_name = "example-template-name"
}

# Or by id:
data "oneuptime_monitor_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `minimum_probe_agreement` (Number) Default minimum number of probes that must agree on a status before the monitor status changes.
- `monitor_description` (String) Default description applied to monitors created from this template.
- `monitor_name` (String) Default name applied to monitors created from this template. Users can override on creation. Leave it blank to name each monitor after the resource it watches.
- `monitor_type` (String) What is the type of monitor created from this template?
- `monitoring_interval` (String) Default monitoring interval for monitors created from this template. A 5-field cron expression, not a label: "*/5 * * * *" is every five minutes.
- `slug` (String) Friendly globally unique name for your object.
- `template_description` (String) Description of the Monitor Template.
- `template_name` (String) Name of the Monitor Template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Default labels applied to monitors created from this template. IDs of `oneuptime_label` resources.
- `monitor_steps` (String) Monitor steps and criteria copied to monitors created from this template.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
