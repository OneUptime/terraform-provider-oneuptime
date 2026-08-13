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
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  status_page_subscriber_notification_template_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `status_page_subscriber_notification_template_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_subscriber_notification_template.example <id>
```
