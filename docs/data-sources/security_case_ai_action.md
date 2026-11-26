---
page_title: "oneuptime_security_case_ai_action Data Source - oneuptime"
subcategory: "Security"
description: |-
  What OneUptime AI did, asked a person about or left to a person on a security case it works, and why.
---

# oneuptime_security_case_ai_action (Data Source)

What OneUptime AI did, asked a person about or left to a person on a security case it works, and why.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case ai action may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_ai_action" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_case_ai_action" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) ID of OneUptime AI's run that did, asked or left this.
- `area` (String) The kind of action the project's setting for it covers: Case, Edge, Accounts, Hosts, Code or Close. Empty for a step no setting covers.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) Unused: OneUptime AI writes every action. The ID of a `oneuptime_user` (see the data source).
- `decided_by_user_id` (String) ID of the person who approved or declined it.
- `detection_rule_id` (String) ID of the detection rule that opened the case: the rule an exception is offered for, and what makes cases alike.
- `hand_off_reason` (String) For a hand-off to people: why OneUptime AI handed the case over.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What OneUptime AI did, or would have done: TakeCase, SetPriority, MarkStepDone, MarkContained, BlockAddress, ContainAccount, ContainHost, DeclareIncident, RequestCodeFix, CloseCase, LiftBlock, OfferRuleException, StepForPeople or HandToPeople.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `mode` (String) What the project's setting for this kind of action said when OneUptime AI acted: Off, Ask or OnItsOwn.
- `reason` (String) Why it is as it is, as a code the case's AI card words.
- `report` (String) For a close: OneUptime AI's closing report, in Markdown - what happened, the impact, what was done, what was fixed and what to watch.
- `security_case_id` (String) ID of the case OneUptime AI worked. The ID of a `oneuptime_security_case`.
- `security_playbook_run_id` (String) ID of the playbook engine's run that carries it out, or waits for its approval.
- `status` (String) What became of it: Doing, Done, Asked (waiting for a person), Declined, Left (OneUptime AI would not), Failed or Undone.
- `step_key` (String) The step of the case's Response Plan it is about: what the step is, and what it is about.
- `target` (String) What it acts on, as the case names it: an address, a user, a host, a route, or the status it sets.
- `undone_by_user_id` (String) ID of the person who undid it.

### Read-Only

- `asked_at` (String) When OneUptime AI asked a person to approve it.
- `created_at` (String) Date and Time when the object was created.
- `decided_at` (String) When a person approved or declined it.
- `done_at` (String) When it was done.
- `evidence` (String) What OneUptime AI read that made it do, ask or leave this: facts, checks and links. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this action belongs to. The ID of a `oneuptime_project`.
- `undone_at` (String) When a person undid it.
- `updated_at` (String) Date and Time when the object was updated.
