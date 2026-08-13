---
page_title: "oneuptime_workspace_notification_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Notification Rule for Third Party Workspaces
---

# oneuptime_workspace_notification_rule (Data Source)

Notification Rule for Third Party Workspaces Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workspace_notification_rule" "by_name" {
  name = "example-workspace_notification_rule"
}

data "oneuptime_workspace_notification_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of the Notification Rule.. Computed.
- `notification_rule` (String) Notification Rules for the Workspace.. Computed.
- `event_type` (String) Event Type for the Workspace like Incident Created, Monitor Status Updated, etc... Computed.
- `workspace_type` (String) Type of Workspace - slack, microsoft teams etc... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
