---
page_title: "oneuptime_automatic_response_action Data Source - oneuptime"
subcategory: "Other"
description: |-
  What each automatic response policy did about an address: challenged or blocked it (or, in its trial run, would have), or a guardrail spared it - and how long the block took.
---

# oneuptime_automatic_response_action (Data Source)

What each automatic response policy did about an address: challenged or blocked it (or, in its trial run, would have), or a guardrail spared it - and how long the block took.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one automatic response action may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_automatic_response_action" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_automatic_response_action" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `block_action` (String) What the policy did, or would have done: Block or Challenge.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `detail` (String) Why the policy left the address alone or could not block it, in a sentence: what protects it, who signed in from it, or what failed.
- `detection_title` (String) The detection the policy acted on, as its case names it.
- `doubt_detail` (String) Who was seen signed in from the address after the policy acted on it.
- `duration` (String) How long the block lasts: the policy's first length, longer for an address blocked before in the last 30 days.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_trial_run` (Boolean) Decided in the policy's trial run: recorded, and nothing done.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `outcome` (String) What the policy did: Challenged or Blocked (in a trial run, what it would have done), or why it left the address alone.
- `policy_name` (String) The policy's name when it acted, kept if the policy is renamed or deleted.
- `prior_block_count` (Number) How many times the address, or a range holding it, was blocked in the 30 days before.
- `security_automatic_response_policy_id` (String) ID of the policy that acted. The ID of a `oneuptime_automatic_response_policy`.
- `security_block_id` (String) ID of the block the policy made. The ID of a `oneuptime_block`.
- `security_case_id` (String) ID of the case whose detection the policy acted on. The ID of a `oneuptime_security_case`.
- `security_playbook_run_id` (String) ID of the run of the playbook engine that made the block. The ID of a `oneuptime_security_playbook_run`.
- `time_to_block_in_seconds` (Number) Seconds from the address's first request to the block taking hold.
- `value` (String) The address the policy acted on: 35.205.11.131.

### Read-Only

- `applied_at` (String) When the block first took hold at a destination that enforces the blocklist: Cloudflare, AWS WAF or a firewall.
- `blocked_at` (String) When the policy's block was made on the blocklist.
- `created_at` (String) Date and Time when the object was created.
- `doubt_seen_at` (String) When a signed-in user was seen from the address after the policy acted on it: the policy lifted its block then (in a trial run, it would have).
- `enforced_at` (String) Where the block took hold, by name, and whether an edge that cannot challenge holds a challenge as a block there. A JSON value: read it with `jsondecode()`.
- `first_request_at` (String) When the address's first request of the detection happened: time to block is measured from here.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `paged_at` (String) For an action stopped at the hourly limit: when the project's security owners were paged about it.
- `project_id` (String) ID of the project the policy acted for. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
