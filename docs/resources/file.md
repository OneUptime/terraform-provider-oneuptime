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
  name = "Example short text"
  file_type = "Example short text"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `file_type` (String) File file_type.

### Optional

- `file` (String) File file.
- `slug` (String) File slug.

### Read-Only

- `id` (String) Unique identifier for the resource.
- `is_public` (Bool) Whether anyone may read the file without signing in. Set by OneUptime: every upload starts private, and a file becomes public only when a record that shows it to everyone, such as a public note or a probe's icon, is published...
- `image_access_token` (String) File image_access_token.

## Import

This resource does not support import: the OneUptime API exposes no read endpoint for it.
