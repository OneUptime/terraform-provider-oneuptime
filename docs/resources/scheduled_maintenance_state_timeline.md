---
page_title: "oneuptime_scheduled_maintenance_state_timeline Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Change state of your maintenance event.
---

# oneuptime_scheduled_maintenance_state_timeline (Resource)

Change state of your maintenance event.

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_state_timeline" "example" {
  scheduled_maintenance_id       = oneuptime_scheduled_maintenance_event.example.id
  scheduled_maintenance_state_id = oneuptime_scheduled_maintenance_state.example.id
}
```

## Schema

### Required

- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance this resource belongs to. The ID of a `oneuptime_scheduled_maintenance_event`.
- `scheduled_maintenance_state_id` (String) ID of the state the maintenance event changes to. The ID of a `oneuptime_scheduled_maintenance_state`.

### Optional

- `ends_at` (String) When did this status change end? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this state change?
- `starts_at` (String) When did this status change? Correct this when the recorded time is wrong - every measurement derived from this timeline is recomputed from the corrected value.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this scheduled maintenance state change.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance state timeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_state_timeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_state_timeline.example <id>
```
