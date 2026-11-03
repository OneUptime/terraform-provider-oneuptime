---
page_title: "oneuptime_service_team_owner Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Add teams as owners to your Service.
---

# oneuptime_service_team_owner (Resource)

Add teams as owners to your Service.

## Example Usage

```terraform
resource "oneuptime_service_team_owner" "example" {
  team_id    = oneuptime_team.example.id
  service_id = oneuptime_service.example.id
}
```

## Schema

### Required

- `service_id` (String) ID of the service the team owns. The ID of a `oneuptime_service`.
- `team_id` (String) ID of the team that owns the service. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service team owner by its service and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_team_owner.example
  id = "<service_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_team_owner.example <service_id>/<team_id>
```

Its own ID works too.
