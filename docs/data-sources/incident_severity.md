---
page_title: "oneuptime_incident_severity Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  How much an incident matters. A new project starts with Critical Incident, Major Incident and Minor Incident: rename them, recolor them or add your own.
---

# oneuptime_incident_severity (Data Source)

How much an incident matters. A new project starts with Critical Incident, Major Incident and Minor Incident: rename them, recolor them or add your own.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident severity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_severity" "example" {
  name = "Example incident severity"
}

# Or by id:
data "oneuptime_incident_severity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Archived severities stay on the incidents and episodes that have them, but can't be picked again. A project keeps at least one severity that is not archived.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this severity ranks among the project's incident severities: 1 is the most severe. A new severity without a number goes to the end of the list. Setting a number places the severity there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When this incident severity was archived. Empty while it is not archived.
- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
