---
page_title: "oneuptime_monitor_label_rule Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Configure rules for automatically attaching labels to monitors when matching monitors are created
---

# oneuptime_monitor_label_rule (Data Source)

Configure rules for automatically attaching labels to monitors when matching monitors are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_label_rule" "example" {
  name = "Example monitor label rule"
}

# Or by id:
data "oneuptime_monitor_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this monitor label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to match any description.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to match any name.
- `name` (String) Name of this monitor label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `labels_to_add` (Set of String) Labels to attach to the monitor when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only trigger for monitors that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
