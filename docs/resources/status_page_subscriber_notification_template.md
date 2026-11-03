---
page_title: "oneuptime_status_page_subscriber_notification_template Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.
---

# oneuptime_status_page_subscriber_notification_template (Resource)

Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.

~> **Renamed:** this resource was called `oneuptime_subscriber_notification_template` before, a name [another resource](./subscriber_notification_template) has now. Rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing status page subscriber notification template:

```terraform
moved {
  from = oneuptime_subscriber_notification_template.example
  to   = oneuptime_status_page_subscriber_notification_template.example
}
```

Where the engine cannot move a resource to another type (Terraform before 1.8), remove it from the state with `terraform state rm` and import it under its new name instead.

## Example Usage

```terraform
resource "oneuptime_status_page_subscriber_notification_template" "example" {
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

Import an existing status page subscriber notification template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_subscriber_notification_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_subscriber_notification_template.example <id>
```
