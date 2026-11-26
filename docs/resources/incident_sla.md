---
page_title: "oneuptime_incident_sla Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Track SLA status and deadlines for incidents
---

# oneuptime_incident_sla (Resource)

Track SLA status and deadlines for incidents

## Example Usage

```terraform
resource "oneuptime_incident_sla" "example" {
  incident_id          = oneuptime_incident.example.id
  incident_sla_rule_id = oneuptime_incident_sla_rule.example.id
}
```

## Schema

### Required

- `incident_id` (String) ID of the incident this SLA record is tracking. The ID of a `oneuptime_incident`.
- `incident_sla_rule_id` (String) ID of the SLA rule that was applied to this incident. The ID of a `oneuptime_incident_sla_rule`.

### Read-Only

- `breach_notification_sent_at` (String) The time when breach notification was sent to incident owners.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_internal_note_reminder_sent_at` (String) The last time an internal note reminder was sent.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_public_note_reminder_sent_at` (String) The last time a public note reminder was sent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resolution_deadline` (String) The deadline by which the incident must be resolved to meet the SLA.
- `resolved_at` (String) The actual time when the incident was resolved.
- `responded_at` (String) The actual time when the incident was acknowledged.
- `response_deadline` (String) The deadline by which the incident must be acknowledged to meet the SLA.
- `sla_started_at` (String) The time when SLA tracking started (usually the incident declaredAt time).
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident sla by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_sla.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_sla.example <id>
```
