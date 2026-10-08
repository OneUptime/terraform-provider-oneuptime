---
page_title: "oneuptime_monitor_group_resource Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Add monitors to your monitor group
---

# oneuptime_monitor_group_resource (Data Source)

Add monitors to your monitor group

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor group resource may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_group_resource" "example" {
  monitor_group_id = oneuptime_monitor_group.example.id
}

# Or by id:
data "oneuptime_monitor_group_resource" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_group_id` (String) ID of your Monitor Group resource where this object belongs. The ID of a `oneuptime_monitor_group`.
- `monitor_id` (String) Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
