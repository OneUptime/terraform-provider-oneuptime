---
page_title: "oneuptime_security_case_step_check Data Source - oneuptime"
subcategory: "Security"
description: |-
  What OneUptime found when it checked a step of a security case's Response Plan itself, and the evidence for it.
---

# oneuptime_security_case_step_check (Data Source)

What OneUptime found when it checked a step of a security case's Response Plan itself, and the evidence for it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case step check may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_step_check" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_case_step_check" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) Unused: OneUptime writes every check. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `security_case_id` (String) ID of the case whose step was checked. The ID of a `oneuptime_security_case`.
- `state` (String) What OneUptime found: Verified, Looks Done (it fits, but OneUptime cannot prove who did it), Not Done, Cannot Tell (no data to check with), Regressed (it was done, and stopped holding) or Not Needed.
- `step_key` (String) The step of the case's Response Plan that was checked: what it is, and what it is about.

### Read-Only

- `checked_at` (String) When OneUptime last checked the step.
- `created_at` (String) Date and Time when the object was created.
- `evidence` (String) Why OneUptime says so: the sentence the step says, what else it saw, links to the traces, events or product, and the figures it worked from. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `next_check_at` (String) When OneUptime next checks the step. Empty when it no longer checks it: the case is closed, or the step left the plan.
- `project_id` (String) ID of the project this check belongs to. The ID of a `oneuptime_project`.
- `state_since` (String) Since when the step has been in its state.
- `updated_at` (String) Date and Time when the object was updated.
