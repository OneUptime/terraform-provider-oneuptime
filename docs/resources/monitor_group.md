---
page_title: "oneuptime_monitor_group Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.
---

# oneuptime_monitor_group (Resource)

Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.

## Example Usage

```terraform
resource "oneuptime_monitor_group" "example" {
  name        = "Example monitor group"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name for this monitor group.

### Optional

- `description` (String) Friendly description that will help you remember.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing monitor group by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_group.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_group.example <id>
```
