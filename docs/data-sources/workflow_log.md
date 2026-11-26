---
page_title: "oneuptime_workflow_log Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Logs of the workflows executed
---

# oneuptime_workflow_log (Data Source)

Logs of the workflows executed

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workflow log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workflow_log" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_workflow_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `can_run_again` (Boolean) True when OneUptime kept what this run started with, so a person can run the workflow again with the same values.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_test_run` (Boolean) True for a run a person started by hand from the Builder while the workflow was turned off. Its steps ran for real; the workflow's trigger stayed paused.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `logs` (String) Logs.
- `run_as` (String) What the run's steps could do in the project: Default, ProjectOwner or CustomRole (runAsCustomRoleName names the role). Set by OneUptime as the run starts.
- `run_as_custom_role_name` (String) The name of the custom role the run's steps acted as, as it was when the run started. Empty unless the run ran as a custom role.
- `run_as_given_by_api_key_id` (String) ID of the project API key whose permissions the run ran with, when it ran as more than Default and a key gave them.
- `run_as_given_by_api_key_name` (String) The name of the project API key whose permissions the run ran with, as it was when the run started. Empty unless a key gave them.
- `run_as_given_by_user_id` (String) ID of the person whose permissions the run ran with, when it ran as more than Default. The ID of a `oneuptime_user` (see the data source).
- `started_by_user_id` (String) ID of the person who started this run by hand from the Builder. Empty for a run its trigger or another workflow started. The ID of a `oneuptime_user` (see the data source).
- `workflow_id` (String) ID of Workflow this logs belong to. The ID of a `oneuptime_workflow`.
- `workflow_status` (String) Status of this workflow.

### Read-Only

- `completed_at` (String) When did this workflow complete.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resume_at` (String) When this workflow run is scheduled to resume after a Sleep step (only set while the run is Waiting).
- `started_at` (String) When did this workflow start.
- `step_trace` (String) Structured per-step record of this run: arguments, return values, the port taken and timing. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
