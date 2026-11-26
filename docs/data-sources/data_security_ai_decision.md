---
page_title: "oneuptime_data_security_ai_decision Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Why OneUptime AI did, or did not, act on a Data Security issue: one record each time it evaluated the issue, with what each kind of work did or why it did nothing - no rule matched, the work is off, no credits, the daily limit, nothing it may do, a proof that failed. Written by OneUptime.
---

# oneuptime_data_security_ai_decision (Data Source)

Why OneUptime AI did, or did not, act on a Data Security issue: one record each time it evaluated the issue, with what each kind of work did or why it did nothing - no rule matched, the work is off, no credits, the daily limit, nothing it may do, a proof that failed. Written by OneUptime.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security ai decision may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_ai_decision" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_security_ai_decision" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_security_issue_id` (String) The issue OneUptime AI evaluated. The ID of a `oneuptime_data_security_issue` (see the data source).
- `data_store_id` (String) The store the issue is on. The ID of a `oneuptime_data_store`.
- `did_act` (Boolean) Whether OneUptime AI did anything for the issue in this evaluation: analyzed it, proposed a change, or made one without asking.
- `did_analyze` (Boolean) Whether OneUptime AI analyzed the issue in this evaluation.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `issue_reopen_count` (Number) How many times the issue had opened again when it was evaluated: OneUptime AI analyzes an issue once each time it opens.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `skip_reason` (String) When OneUptime AI did nothing: RulesUnreadable, NoRuleMatched, StoreDailyLimitReached, ProjectDailyLimitReached, NothingToDo (each kind of work says why in the entries) or EvaluationFailed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `entries` (String) One entry per kind of work - Analyze, Contain, Fix, Judge - with a reason code saying what it did or why it did nothing, and the rule or proposal it was about. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this decision belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
