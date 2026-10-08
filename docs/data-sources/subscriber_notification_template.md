---
page_title: "oneuptime_subscriber_notification_template Data Source - oneuptime"
subcategory: "Other"
description: |-
  Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.
---

# oneuptime_subscriber_notification_template (Data Source)

Links subscriber notification templates to specific status pages. This allows you to use different notification templates for different status pages.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one subscriber notification template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_subscriber_notification_template" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_subscriber_notification_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `status_page_id` (String) ID of the Status Page this template is linked to. The ID of a `oneuptime_status_page`.
- `status_page_subscriber_notification_template_id` (String) ID of the notification template linked to this status page. The ID of a `oneuptime_subscriber_notification_template`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
