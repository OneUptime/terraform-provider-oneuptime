---
page_title: "oneuptime_file Resource - oneuptime"
subcategory: "Organization"
description: |-
  BLOB or File storage
---

# oneuptime_file (Resource)

BLOB or File storage

## Example Usage

```terraform
resource "oneuptime_file" "example" {
  name      = "Example file"
  file_type = "Example short text"
}
```

## Schema

### Required

- `file_type` (String)
- `name` (String) Any friendly name of this object.

### Optional

- `file` (String)
- `slug` (String)

### Read-Only

- `id` (String) Unique identifier for the resource.
- `image_access_token` (String)
- `is_public` (Boolean) Whether anyone may read the file without signing in. Set by OneUptime: every upload starts private, and a file becomes public only when a record that shows it to everyone, such as a public note or a probe's icon, is published.

## Import

This resource does not support import: the OneUptime API exposes no read endpoint for it.
