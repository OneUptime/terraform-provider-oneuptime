---
page_title: "oneuptime_subscriber_notification_template Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage custom notification templates for status page subscribers. These templates can be used to customize the notifications sent to subscribers via Email, SMS, Slack, Microsoft Teams, and Webhooks.
---

# oneuptime_subscriber_notification_template (Resource)

Manage custom notification templates for status page subscribers. These templates can be used to customize the notifications sent to subscribers via Email, SMS, Slack, Microsoft Teams, and Webhooks.

## Example Usage

```terraform
resource "oneuptime_subscriber_notification_template" "example" {
  template_name       = "Example short text"
  event_type          = "Subscriber Subscription Confirmation"
  notification_method = "Email"
  template_body       = "<div><h1>Title</h1><p>Content</p></div>"
}
```

## Schema

### Required

- `event_type` (String) The type of event this template is for (e.g., Incident Created, Announcement Created). Allowed values: `Subscriber Subscription Confirmation`, `Subscriber Subscribed`, `Subscriber Manage Subscription`, `Subscriber Incident Created`, `Subscriber Incident State Changed`, `Subscriber Incident Note Created`, `Subscriber Incident Note Updated`, `Subscriber Incident Postmortem Published`, `Subscriber Episode Created`, `Subscriber Episode State Changed`, `Subscriber Episode Note Created`, `Subscriber Episode Note Updated`, `Subscriber Announcement Created`, `Subscriber Announcement Updated`, `Subscriber Scheduled Maintenance Created`, `Subscriber Scheduled Maintenance State Changed`, `Subscriber Scheduled Maintenance Note Created`, `Subscriber Scheduled Maintenance Note Updated`, `Subscriber Report`.
- `notification_method` (String) The notification method this template is for (Email, SMS, Slack, Microsoft Teams, Webhook). Allowed values: `Email`, `SMS`, `Slack`, `Microsoft Teams`, `Webhook`.
- `template_body` (String) The template body content. For Email: HTML template. For SMS: Plain text. For Slack/Teams: Markdown.
- `template_name` (String) A friendly name for this notification template.

### Optional

- `email_subject` (String) The subject line for email notifications. Only used for Email notification method.
- `template_description` (String) A description for this notification template.

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
