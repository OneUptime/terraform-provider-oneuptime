---
page_title: "oneuptime_cloud_resource_team_owner Data Source - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your cloud resources.
---

# oneuptime_cloud_resource_team_owner (Data Source)

Add teams as owners to your cloud resources.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud resource team owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_resource_team_owner" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_cloud_resource_team_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `cloud_resource_id` (String) ID of your OneUptime Cloud Resource in which this object belongs. The ID of a `oneuptime_cloud_resource`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
