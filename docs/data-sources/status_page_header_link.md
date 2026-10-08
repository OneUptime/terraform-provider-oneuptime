---
page_title: "oneuptime_status_page_header_link Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage header links on your status page
---

# oneuptime_status_page_header_link (Data Source)

Manage header links on your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page header link may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_header_link" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_status_page_header_link" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `link` (String) URL to a website or any other resource on the internet.
- `order` (Number) Where this link appears among the status page's header links, lowest number first. A new link is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `title` (String) Title of this resource.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
