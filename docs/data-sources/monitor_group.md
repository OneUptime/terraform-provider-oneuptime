---
page_title: "oneuptime_monitor_group Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.
---

# oneuptime_monitor_group (Data Source)

Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor group may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_group" "example" {
  name = "Example monitor group"
}

# Or by id:
data "oneuptime_monitor_group" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name for this monitor group.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
