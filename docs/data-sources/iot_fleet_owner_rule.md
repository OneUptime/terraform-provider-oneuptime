---
page_title: "oneuptime_iot_fleet_owner_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching IoT fleets are created
---

# oneuptime_iot_fleet_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching IoT fleets are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one iot fleet owner rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_io_t_fleet_owner_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_iot_fleet_owner_rule" "example" {
  name = "Example iot fleet owner rule"
}

# Or by id:
data "oneuptime_iot_fleet_owner_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this IoT fleet owner rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `iot_fleet_description_pattern` (String) Regex (case-insensitive) matched against the IoT fleet description. Leave empty to match any description.
- `iot_fleet_name_pattern` (String) Regex (case-insensitive) matched against the IoT fleet name. Leave empty to match any name.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this IoT fleet owner rule.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `iot_fleet_labels` (Set of String) Only trigger for IoT fleets that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `owner_teams` (Set of String) Teams to add as owners on the IoT fleet when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the IoT fleet when this rule matches. IDs of `oneuptime_user` records.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
