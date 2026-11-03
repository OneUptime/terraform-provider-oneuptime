---
page_title: "oneuptime_alert_team_owner Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Add teams as owners to your alerts.
---

# oneuptime_alert_team_owner (Resource)

Add teams as owners to your alerts.

## Example Usage

```terraform
resource "oneuptime_alert_team_owner" "example" {
  team_id  = oneuptime_team.example.id
  alert_id = oneuptime_alert.example.id
}
```

## Schema

### Required

- `alert_id` (String) ID of the alert the team owns. The ID of a `oneuptime_alert`.
- `team_id` (String) ID of the team that owns the alert. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert team owner by its alert and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_team_owner.example
  id = "<alert_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_team_owner.example <alert_id>/<team_id>
```

Its own ID works too.
