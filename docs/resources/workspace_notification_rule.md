---
page_title: "oneuptime_workspace_notification_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Notification Rule for Third Party Workspaces
---

# oneuptime_workspace_notification_rule (Resource)

Notification Rule for Third Party Workspaces

## Example Usage

```terraform
resource "oneuptime_workspace_notification_rule" "example" {
  name = "This is an example of longer text content that might be stored in this field."
  event_type = "Example short text"
  workspace_type = "This is an example of longer text content that might be stored in this field."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of the Notification Rule..
- `event_type` (String) Event Type for the Workspace like Incident Created, Monitor Status Updated, etc...
- `workspace_type` (String) Type of Workspace - slack, microsoft teams etc...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of the Notification Rule..
- `notification_rule` (String) Notification Rules for the Workspace..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_workspace_notification_rule.example <id>
```
