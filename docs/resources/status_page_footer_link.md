---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_status_page_footer_link Resource - oneuptime"
subcategory: ""
description: |-
  Status page footer link resource
---

# oneuptime_status_page_footer_link (Resource)

Status page footer link resource

## Example Usage

```terraform
resource "oneuptime_status_page_footer_link" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  title = "example-title"
  link = "example-link"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `title` (String) Title. Required.
- `link` (String) Link. Required.
- `order` (Number) Order. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_footer_link.example <id>
```
