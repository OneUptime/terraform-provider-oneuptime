---
page_title: "oneuptime_status_page_footer_link Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage footer links on your status page
---

# oneuptime_status_page_footer_link (Resource)

Manage footer links on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_footer_link" "example" {
  status_page_id = oneuptime_status_page.example.id
  title          = "Example short text"
  link           = "https://short.url/abc123"
}
```

## Schema

### Required

- `link` (String) URL to a website or any other resource on the internet.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `title` (String) Title of this resource.

### Optional

- `order` (Number) Where this link appears among the status page's footer links, lowest number first. A new link is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page footer link by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_footer_link.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_footer_link.example <id>
```
