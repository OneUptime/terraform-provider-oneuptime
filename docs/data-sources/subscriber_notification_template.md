---
page_title: "oneuptime_subscriber_notification_template Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage custom notification templates for status page subscribers. These templates can be used to customize the notifications sent to subscribers via Email, SMS, Slack, Microsoft Teams, and Webhooks.
---

# oneuptime_subscriber_notification_template (Data Source)

Manage custom notification templates for status page subscribers. These templates can be used to customize the notifications sent to subscribers via Email, SMS, Slack, Microsoft Teams, and Webhooks.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one subscriber notification template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_subscriber_notification_template" "example" {
  template_name = "example-template-name"
}

# Or by id:
data "oneuptime_subscriber_notification_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `email_subject` (String) The subject line for email notifications. Only used for Email notification method.
- `event_type` (String) The type of event this template is for (e.g., Incident Created, Announcement Created).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `notification_method` (String) The notification method this template is for (Email, SMS, Slack, Microsoft Teams, Webhook).
- `template_body` (String) The template body content. For Email: HTML template. For SMS: Plain text. For Slack/Teams: Markdown.
- `template_description` (String) A description for this notification template.
- `template_name` (String) A friendly name for this notification template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
