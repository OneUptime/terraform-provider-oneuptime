---
page_title: "oneuptime_status_page_announcement_template Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcement templates for your status page
---

# oneuptime_status_page_announcement_template (Resource)

Manage announcement templates for your status page

## Example Usage

```terraform
resource "oneuptime_status_page_announcement_template" "example" {
  template_name = "Example short text"
  title         = "Example short text"
  description   = "Managed by Terraform"
}
```

## Schema

### Required

- `description` (String) Text of the announcement. This is in Markdown.
- `template_name` (String) Name of the announcement template.
- `title` (String) Title of the announcement.

### Optional

- `monitors` (Set of String) List of monitors affected by this announcement template. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about announcements created from this template? Defaults to `true`.
- `status_pages` (Set of String) Status Pages to show this announcement on. IDs of `oneuptime_status_page` resources.
- `template_description` (String) Description of the announcement template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page announcement template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_announcement_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_announcement_template.example <id>
```
