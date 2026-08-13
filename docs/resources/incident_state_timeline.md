---
page_title: "oneuptime_incident_state_timeline Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Change state of the incidents (Created to Acknowledged for example)
---

# oneuptime_incident_state_timeline (Resource)

Change state of the incidents (Created to Acknowledged for example)

## Example Usage

```terraform
resource "oneuptime_incident_state_timeline" "example" {
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
  incident_state_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `incident_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_state_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about this state change?..
- `root_cause` (String) What is the root cause of this status change?..
- `ends_at` (String) A date time object..
- `starts_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this incident state change..
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons..
- `is_owner_notified` (Bool) Are owners notified of state change?..
- `state_change_log` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident State Timeline], Update: [No access - you don't have permission for this operation].

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_state_timeline.example <id>
```
