---
page_title: "oneuptime_form_submission Data Source - oneuptime"
subcategory: "Forms"
description: |-
  The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident, alert or maintenance event each one created.
---

# oneuptime_form_submission (Data Source)

The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident, alert or maintenance event each one created.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one form submission may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_form_submission" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_form_submission" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert this submission created. Empty for a form that creates incidents or schedules maintenance, and once that alert is deleted. The ID of a `oneuptime_alert`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `form_id` (String) ID of the form this submission was made through. The ID of a `oneuptime_form`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this submission created. Empty for a form that creates alerts or schedules maintenance, and once that incident is deleted. The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `scheduled_maintenance_id` (String) ID of the maintenance event this submission created. Empty for a form that creates incidents or alerts, and once that event is deleted. The ID of a `oneuptime_scheduled_maintenance_event`.
- `submitter_name` (String) The name the submitter gave, as they typed it. Empty when the form does not ask for it, or lets people leave it out and they did.
- `target_type` (String) What the submission created: Incident, Alert, or ScheduledMaintenance (a maintenance event).

### Read-Only

- `answers` (String) Every question the submitter answered, in the form's order: the question's id (fieldId), its label when the form was submitted, the stored value, and the value as a person reads it (displayValue). A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `submitter_email` (String) The email address the submitter gave. It is not verified, and nothing is sent to it. Empty when the form does not ask for it, or lets people leave it out and they did.
- `updated_at` (String) Date and Time when the object was updated.
