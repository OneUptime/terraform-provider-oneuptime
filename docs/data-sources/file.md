---
page_title: "oneuptime_file Data Source - oneuptime"
subcategory: "Organization"
description: |-
  BLOB or File storage
---

# oneuptime_file (Data Source)

BLOB or File storage

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one file may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_file" "example" {
  name = "Example file"
}

# Or by id:
data "oneuptime_file" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `file` (String) The file's content.
- `file_type` (String) The file's media type, such as image/png.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `image_access_token` (String) Addresses the file in a link (/file/image/access-token/<token>) without its ID. Opening it still takes the right to see the file. Set by OneUptime.
- `is_public` (Boolean) Whether anyone may read the file without signing in. Set by OneUptime: every upload starts private, and a file becomes public only when a record that shows it to everyone, such as a public note or a probe's icon, is published.
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.
