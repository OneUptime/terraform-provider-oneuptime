---
page_title: "oneuptime_status_page_team_owner Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Add teams as owners to your Status Page.
---

# oneuptime_status_page_team_owner (Resource)

Add teams as owners to your Status Page.

## Example Usage

```terraform
resource "oneuptime_status_page_team_owner" "example" {
  team_id        = oneuptime_team.example.id
  status_page_id = oneuptime_status_page.example.id
}
```

## Schema

### Required

- `status_page_id` (String) ID of the status page the team owns. The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of the team that owns the status page. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page team owner by its status page and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_team_owner.example
  id = "<status_page_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_team_owner.example <status_page_id>/<team_id>
```

Its own ID works too.
