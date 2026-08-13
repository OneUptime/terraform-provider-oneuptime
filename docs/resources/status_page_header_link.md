---
page_title: "oneuptime_status_page_header_link Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage header links on your status page
---

# oneuptime_status_page_header_link (Resource)

Manage header links on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_header_link" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  title = "Example short text"
  link = "https://short.url/abc123"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `title` (String) Title of this resource..
- `link` (String) URL to a website or any other resource on the internet..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this resource..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_header_link.example <id>
```
