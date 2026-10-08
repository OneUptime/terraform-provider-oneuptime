---
page_title: "oneuptime_label Resource - oneuptime"
subcategory: "Organization"
description: |-
  Organize resources for your project by using labels / tags.
---

# oneuptime_label (Resource)

Organize resources for your project by using labels / tags.

## Example Usage

```terraform
resource "oneuptime_label" "example" {
  name        = "Example label"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing label by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_label.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_label.example <id>
```
