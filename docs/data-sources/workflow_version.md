---
page_title: "oneuptime_workflow_version Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Every saved version of a workflow's graph: who saved it, when, and what it changed, so an earlier one can be looked at and restored.
---

# oneuptime_workflow_version (Data Source)

Every saved version of a workflow's graph: who saved it, when, and what it changed, so an earlier one can be looked at and restored.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workflow version may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workflow_version" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_workflow_version" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `change_count` (Number) How many changes this version made, counted before its lines were cut to 25.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who saved this version. The ID of a `oneuptime_user` (see the data source).
- `graph_fingerprint` (String) SHA-256 of the graph without what only the builder's canvas keeps (which card is selected, how big each one is). Two versions with the same graph have the same fingerprint.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `restored_from_version_number` (Number) For a Restored version: the number of the earlier version it put back.
- `saved_through` (String) AskAi when Ask AI saved this version for the person in Saved by User: each change Ask AI makes is a version of its own. Empty when the person saved it themselves, or the API did.
- `version_number` (Number) The version's number: one more than the version before it. Only the newest 100 versions of a workflow are kept.
- `version_type` (String) What made this version: Created (the workflow as it was created), Changed (changes saved in the builder, by Ask AI or through the API), Restored (an earlier version put back) or Earlier (the workflow as it was when version history began).
- `workflow_id` (String) ID of the workflow this is a version of. The ID of a `oneuptime_workflow`.

### Read-Only

- `changes` (String) What this version changed from the version before it, one plain line per change: steps added, removed or reconfigured, connections made or removed, steps moved. A setting is named, never quoted. At most 25 lines. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `graph` (String) The workflow's graph as this version saved it: its steps, where each one sits, and their connections. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
