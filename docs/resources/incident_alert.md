---
page_title: "oneuptime_incident_alert Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.
---

# oneuptime_incident_alert (Resource)

Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.

## Example Usage

```terraform
resource "oneuptime_incident_alert" "example" {
  incident_id = oneuptime_incident.example.id
  alert_id    = oneuptime_alert.example.id
}
```

## Schema

### Required

- `alert_id` (String) ID of the Alert that is linked to the incident. The ID of a `oneuptime_alert`.
- `incident_id` (String) ID of the Incident this alert is linked to. The ID of a `oneuptime_incident`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the User who linked the alert to the incident (if it was linked by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident alert by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_alert.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_alert.example <id>
```
