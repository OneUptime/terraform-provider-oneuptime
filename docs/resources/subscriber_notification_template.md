---
page_title: "oneuptime_subscriber_notification_template Resource - oneuptime"
subcategory: "Other"
description: |-
  Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.
---

# oneuptime_subscriber_notification_template (Resource)

Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.

## Example Usage

```terraform
resource "oneuptime_subscriber_notification_template" "example" {
  status_page_id                                  = oneuptime_status_page.example.id
  status_page_subscriber_notification_template_id = oneuptime_subscriber_notification_template.example.id
}
```

## Schema

### Required

- `status_page_id` (String) ID of the Status Page this template is linked to. The ID of a `oneuptime_status_page`.
- `status_page_subscriber_notification_template_id` (String) ID of the notification template linked to this status page. The ID of a `oneuptime_subscriber_notification_template`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing subscriber notification template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_subscriber_notification_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_subscriber_notification_template.example <id>
```
