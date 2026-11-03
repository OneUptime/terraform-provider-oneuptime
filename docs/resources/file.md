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
  file_type = "audio/aac"
}
```

## Schema

### Required

- `file_type` (String) The file's media type, such as image/png. Allowed values: `audio/aac`, `image/png`, `image/jpeg`, `image/svg+xml`, `image/gif`, `image/webp`, `image/x-icon`, `image/bmp`, `image/tiff`, `image/avif`, `image/heic`, `application/pdf`, `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`, `text/plain`, `text/markdown`, `text/csv`, `application/rtf`, `application/vnd.oasis.opendocument.text`, `application/json`, `application/zip`, `application/vnd.ms-excel`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`, `application/vnd.oasis.opendocument.spreadsheet`, `application/vnd.ms-powerpoint`, `application/vnd.openxmlformats-officedocument.presentationml.presentation`, `application/vnd.oasis.opendocument.presentation`, `application/vnd.tcpdump.pcap`.
- `name` (String) Any friendly name of this object.

### Optional

- `file` (String) The file's content.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `id` (String) Unique identifier for the resource.
- `image_access_token` (String) Addresses the file in a link (/file/image/access-token/<token>) without its ID. Opening it still takes the right to see the file. Set by OneUptime.
- `is_public` (Boolean) Whether anyone may read the file without signing in. Set by OneUptime: every upload starts private, and a file becomes public only when a record that shows it to everyone, such as a public note or a probe's icon, is published.

## Import

This resource does not support import: the OneUptime API exposes no read endpoint for it.
