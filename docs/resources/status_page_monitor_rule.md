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
  status_page_id = oneuptime_status_page.example.id
  name           = "Example status page monitor rule"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this status page monitor rule.
- `status_page_id` (String) ID of the status page this rule adds matching monitors to. The ID of a `oneuptime_status_page`.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this status page monitor rule.
- `is_enabled` (Boolean) Whether this rule is enabled. A disabled rule removes the monitors it had added. Defaults to `true`.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter.
- `monitor_labels` (Set of String) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded on the resources this rule adds. Defaults to `true`.
- `show_status_history_chart` (Boolean) Show a 90 day uptime history on the resources this rule adds to the status page. Defaults to `true`.
- `show_uptime_percent` (Boolean) Show uptime percent on the resources this rule adds to the status page. Defaults to `true`.
- `status_page_group_id` (String) ID of the group that matched monitors are added to. Empty means ungrouped. The ID of a `oneuptime_status_page_group`.
- `uptime_percent_precision` (String) Precision of the uptime percent shown on the resources this rule adds.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page monitor rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_monitor_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_monitor_rule.example <id>
```
