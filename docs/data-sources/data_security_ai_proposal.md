---
page_title: "oneuptime_data_security_ai_proposal Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  What OneUptime AI proposes for a Data Security issue, or for who owns the project's stores: what will change, why, how it will be proven and how to undo it, and whether it waits for a person, was made, dismissed or undone. Written by OneUptime; approved, dismissed and undone through Data Security's own routes, with the person's own permission for the change.
---

# oneuptime_data_security_ai_proposal (Data Source)

What OneUptime AI proposes for a Data Security issue, or for who owns the project's stores: what will change, why, how it will be proven and how to undo it, and whether it waits for a person, was made, dismissed or undone. Written by OneUptime; approved, dismissed and undone through Data Security's own routes, with the person's own permission for the change.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security ai proposal may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_ai_proposal" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_security_ai_proposal" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) The AI run that made the proposal, when one did.
- `approved_by_user_id` (String) The person who approved it. The ID of a `oneuptime_user` (see the data source).
- `ask_reason` (String) Why a change OneUptime AI could have made on its own waits for a person: the work asks first, the change always asks, a rule must name it, it could not prove it is safe, two people approve it, the hour's circuit breaker, or the repository's cap on open AI pull requests.
- `auto_remediation_rule_id` (String) The AI rule that allowed it; empty when the project's defaults did.
- `auto_remediation_rule_name` (String) What the AI rule that allowed it was called when it was proposed; empty when the project's defaults did.
- `autonomy` (String) What the project's AI rules (or its defaults) allowed for this kind of work when it was proposed: Suggest or WithoutAsking.
- `change_markdown` (String) A preview of the change a person can read before approving it: each rule it makes, and where.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_security_issue_group_id` (String) The group of issues whose cause the proposal fixes once (#10370): one pull request, or one approval for the stores made in the console. Such a proposal names no one issue.
- `data_security_issue_id` (String) The issue the proposal is for. Empty for an owner rule, which is for no one issue, and for the owners of a store that has no issue (#10374). The ID of a `oneuptime_data_security_issue` (see the data source).
- `data_store_id` (String) The store the proposal is for: the issue's store, or the store whose owners it proposes. Empty for an owner rule, which is for the project's stores (#10374). The ID of a `oneuptime_data_store`.
- `dismiss_reason` (String) Why the person dismissed it.
- `dismissed_by_user_id` (String) The person who dismissed it. The ID of a `oneuptime_user` (see the data source).
- `failure_message` (String) What stopped the change, in the words of what refused it. Nothing changed.
- `first_approved_by_user_id` (String) The person who gave the first of two approvals. The ID of a `oneuptime_user` (see the data source).
- `first_approved_plan_fingerprint` (String) The plan the first approval was given to: a second approval runs only that plan.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What OneUptime AI proposes: ScrubTelemetry, the scrub rules that hide a telemetry issue's data from now on.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `plan_fingerprint` (String) A fingerprint of the plan: an approval names the one its person was shown, and a plan that changed since is not approved.
- `plan_key` (String) What makes two proposals the same change: OneUptime AI never proposes one for an issue again, unless the last one expired.
- `proof_markdown` (String) How the change is proven to have worked: the issue's own check, run again - never OneUptime AI's word.
- `reason_markdown` (String) Why OneUptime AI proposes it, with the facts it rests on: the check, the data types and where they were found - never a value.
- `risk_score` (Number) The issue's risk score when OneUptime AI proposed it: what waits for a person is listed riskiest first.
- `safety` (String) What kind of change it is, which decides what may run without asking: MaskOnly (hides data in OneUptime, undone by deleting what it made), OpensPullRequest, RecordsJudgement, ChangesAccess, Deletes, Pages, ChangesWhoIsTold or ChangesCloudSettings. Deleting data, paging people, changing who is told about a store and changing a cloud store's settings always ask.
- `status` (String) Where it stands: Proposed (waiting for a person), Approved (being made), Applied, AppliedWithoutAsking, Dismissed, Expired (no longer needed), Undone or Failed.
- `title` (String) What OneUptime AI proposes, in a line.
- `undo_markdown` (String) What Undo does, and what it cannot take back.
- `undone_by_user_id` (String) The person who undid the change. The ID of a `oneuptime_user` (see the data source).
- `work` (String) The kind of work it is, whose autonomy the project's AI rules set: Analyze, Contain, Fix or Judge.

### Read-Only

- `applied_at` (String) When the change was made: after approval, or without asking.
- `approved_at` (String) When it was approved.
- `created_at` (String) Date and Time when the object was created.
- `dismissed_at` (String) When it was dismissed.
- `expired_at` (String) When it stopped being needed: its issue resolved, or was accepted or set aside, before anyone acted on it.
- `failed_at` (String) When making the change failed.
- `first_approved_at` (String) When the first of two approvals was given.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `plan` (String) The change, as its action reads it to make it: for ScrubTelemetry, each finding's data type and place. Names and ids only. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this proposal belongs to. The ID of a `oneuptime_project`.
- `result` (String) What making it made, for Undo to take back: for ScrubTelemetry, the scrub rules it created. A JSON value: read it with `jsondecode()`.
- `undone_at` (String) When the change was undone.
- `updated_at` (String) Date and Time when the object was updated.
- `verified_at` (String) When the issue's own check passed after the change was made: the proof that it worked.
