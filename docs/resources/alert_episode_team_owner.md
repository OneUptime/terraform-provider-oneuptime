---
page_title: "oneuptime_alert_episode_team_owner Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Add teams as owners to your alert episodes.
---

# oneuptime_alert_episode_team_owner (Resource)

Add teams as owners to your alert episodes.

## Example Usage

```terraform
resource "oneuptime_alert_episode_team_owner" "example" {
  team_id          = oneuptime_team.example.id
  alert_episode_id = oneuptime_alert_episode.example.id
}
```

## Schema

### Required

- `alert_episode_id` (String) ID of your OneUptime Alert Episode in which this object belongs. The ID of a `oneuptime_alert_episode`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Optional

- `is_owner_notified` (Boolean) Are owners notified of this resource ownership? Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_team_owner.example <id>
```
