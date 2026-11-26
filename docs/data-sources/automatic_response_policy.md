---
page_title: "oneuptime_automatic_response_policy Data Source - oneuptime"
subcategory: "Other"
description: |-
  Policies that challenge or block the address a security detection is about as the detection is made, under guardrails: never a protected or own address, never one a signed-in user used, longer blocks for repeat offenders, at most 50 an hour, and a trial run first.
---

# oneuptime_automatic_response_policy (Data Source)

Policies that challenge or block the address a security detection is about as the detection is made, under guardrails: never a protected or own address, never one a signed-in user used, longer blocks for repeat offenders, at most 50 an hour, and a trial run first.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one automatic response policy may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_automatic_response_policy" "example" {
  name = "Example automatic response policy"
}

# Or by id:
data "oneuptime_automatic_response_policy" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `block_action` (String) What the policy does to the address: Block, or Challenge (people pass a challenge, scripts do not). Where an edge cannot challenge, a challenge is enforced there as a block for the same time.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who made the policy. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the policy does, in a sentence.
- `duration` (String) How long the policy's first block of an address lasts: 1 hour, 24 hours or 7 days. An address blocked before in the last 30 days is blocked for longer.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) Which policy it is: one of the ready ones (Scanners, Password Guessing, Known-Bad Addresses), or Custom for one the project made.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `mode` (String) Off (does nothing), Trial Run (records what it would have done, for a week, and does nothing else) or On (acts).
- `mode_changed_by_user_id` (String) ID of the person who last turned the policy on, off or into its trial run. The ID of a `oneuptime_user` (see the data source).
- `name` (String) The policy's name: Scanners.
- `trial_skipped_by_user_id` (String) ID of the person who turned the policy on before its trial run was up. The ID of a `oneuptime_user` (see the data source).
- `turned_on_by_user_id` (String) ID of the person who last turned the policy on. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `match_conditions` (String) What a detection must be for a project's own policy to act on it: its ATT&CK techniques or detection rules, whether threat intel must list the address, and whether every request so far must have been turned away. A ready policy's are its own. A JSON value: read it with `jsondecode()`.
- `mode_changed_at` (String) When the policy was last turned on, off or into its trial run.
- `project_id` (String) ID of the project the policy acts for. The ID of a `oneuptime_project`.
- `trial_ends_at` (String) When the policy's trial run has run its week, and it asks to be turned on.
- `trial_finished_at` (String) When the policy last finished a whole week of trial run. Empty while it has never finished one.
- `trial_skipped_at` (String) When the policy was turned on before a week of trial run, skipping the rest of it.
- `trial_started_at` (String) When the policy's trial run started.
- `turned_on_at` (String) When the policy was last turned on: its blocks run without approval from here.
- `updated_at` (String) Date and Time when the object was updated.
