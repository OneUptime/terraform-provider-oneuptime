---
page_title: "oneuptime_webhook_log Data Source - oneuptime"
subcategory: "Notification Logs"
description: |-
  Logs of all the outbound Webhook requests sent for this project.
---

# oneuptime_webhook_log (Data Source)

Logs of all the outbound Webhook requests sent for this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one webhook log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_webhook_log" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_webhook_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of Alert associated with this request (if any). The ID of a `oneuptime_alert`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of Incident associated with this request (if any). The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_id` (String) ID of Monitor associated with this request (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Policy Escalation Rule associated with this request (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Policy associated with this request (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Schedule associated with this request (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `request_body` (String) JSON body that was POSTed to the webhook URL.
- `response_body` (String) Response body returned by the webhook endpoint (truncated).
- `response_status_code` (Number) HTTP status code returned by the webhook endpoint.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this request (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Status of the Webhook request.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this request (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this request (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this request (if any). The ID of a `oneuptime_team`.
- `user_id` (String) ID of User who initiated this request (if any). The ID of a `oneuptime_user` (see the data source).
- `webhook_url` (String) Where the request was sent: the address's scheme, host and port. The rest of the address is left out, as it can hold the webhook's secret.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
