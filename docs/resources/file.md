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
- `is_public` (Bool) File is_public.
- `image_access_token` (String) File image_access_token.

### Read-Only

- `id` (String) Unique identifier for the resource.

## Import

This resource does not support import: the OneUptime API exposes no read endpoint for it.
