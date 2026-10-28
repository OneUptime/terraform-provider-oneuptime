---
page_title: "oneuptime_dashboard_playlist Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Dashboards shown one after another on a screen, such as a NOC wall, each for the same time.
---

# oneuptime_dashboard_playlist (Resource)

Dashboards shown one after another on a screen, such as a NOC wall, each for the same time.

## Example Usage

```terraform
resource "oneuptime_dashboard_playlist" "example" {
  name = "Example dashboard playlist"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `dashboards` (String) The dashboards shown, in order: a list like [{"dashboardId": "<id>"}]. One may appear more than once. A JSON value: write it with `jsonencode()`.
- `interval_in_seconds` (Number) How long each dashboard is shown before the next one, in seconds: 10 to 86400. Defaults to `60`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dashboard playlist by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dashboard_playlist.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dashboard_playlist.example <id>
```
