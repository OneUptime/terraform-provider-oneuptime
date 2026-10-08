---
page_title: "oneuptime_iot_fleet_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to IoT fleets when matching IoT fleets are created
---

# oneuptime_iot_fleet_label_rule (Data Source)

Configure rules for automatically attaching labels to IoT fleets when matching IoT fleets are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one iot fleet label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_io_t_fleet_label_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_iot_fleet_label_rule" "example" {
  name = "Example iot fleet label rule"
}

# Or by id:
data "oneuptime_iot_fleet_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this IoT fleet label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `iot_fleet_description_pattern` (String) Regex (case-insensitive) matched against the IoT fleet description. Leave empty to match any description.
- `iot_fleet_name_pattern` (String) Regex (case-insensitive) matched against the IoT fleet name. Leave empty to match any name.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this IoT fleet label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `iot_fleet_labels` (Set of String) Only trigger for IoT fleets that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `labels_to_add` (Set of String) Labels to attach to the IoT fleet when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
