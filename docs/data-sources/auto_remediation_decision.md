---
page_title: "oneuptime_auto_remediation_decision Data Source - oneuptime"
subcategory: "AI"
description: |-
  Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal.
---

# oneuptime_auto_remediation_decision (Data Source)

Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one auto remediation decision may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_auto_remediation_decision" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_auto_remediation_decision" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert this decision is about. The ID of a `oneuptime_alert`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this decision is about. The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `stage` (String) WaitingForInvestigation while remediation waits for the signal's AI investigation to finish, Evaluated once every fix path was evaluated.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `entries` (String) One entry per fix path - a Kubernetes cluster, an infrastructure resource, an AI rule, or the project - with a reason code saying what it did or why it did nothing. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
