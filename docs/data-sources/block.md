---
page_title: "oneuptime_block Data Source - oneuptime"
subcategory: "Other"
description: |-
  The IP addresses, ranges, domains, URLs and files the project blocks: why, until when, and where each block is enforced.
---

# oneuptime_block (Data Source)

The IP addresses, ranges, domains, URLs and files the project blocks: why, until when, and where each block is enforced.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one block may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_block" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_block" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) Block turns it away. Challenge, for an address at the edge, asks for a challenge a person passes and a script or an API client does not. BlockAndQuarantine, for a file, also quarantines the copies an endpoint product finds.
- `challenge_falls_back_to_block` (Boolean) For a challenge: whether an edge that cannot challenge holds it as a block for the same time. Set for an automatic response policy's challenges.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who blocked it. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `indicator_type` (String) What the block holds: IpAddress, IpRange, Domain, Url, or a file's FileSha256, FileSha1 or FileMd5 hash. Worked out from the value.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `lift_reason` (String) Why the block was lifted, given with the lift.
- `lifted_by_user_id` (String) ID of the person who lifted the block. The ID of a `oneuptime_user` (see the data source).
- `reach` (String) Where it is enforced. Edge, for an address: at the edge in front of the app. Everywhere: on every machine too, through Microsoft Defender or CrowdStrike, which a domain, a URL and a file always are.
- `reason` (String) Why the address is blocked, for whoever reads the list later.
- `security_automatic_response_policy_id` (String) ID of the automatic response policy that made the block. The ID of a `oneuptime_automatic_response_policy`.
- `security_case_id` (String) ID of the case the block was made for. The block and its lifting go on that case's timeline. The ID of a `oneuptime_security_case`.
- `security_playbook_run_id` (String) ID of the playbook run whose step made the block. The ID of a `oneuptime_security_playbook_run`.
- `source` (String) Who or what made the block: a Person, a Playbook, an AutomaticResponse policy, OneUptimeAi or an Api key.
- `state` (String) Active until it expires or someone lifts it. Set it to Lifted to lift the block; Expired is set by OneUptime when its time is up.
- `value` (String) What is blocked: an IPv4 or IPv6 address, a range such as 203.0.113.0/24, a domain, an http or https URL, or a file's MD5, SHA-1 or SHA-256 hash. Kept one way: a range as its first address and its length, an IPv6 address compressed in lowercase, a domain lower case without a trailing dot, a URL as a browser sends it, a hash in lower-case hex.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `expires_at` (String) When the block ends by itself: 24 hours after it is made unless the request says otherwise. Null keeps it until someone lifts it, which takes a Security Admin.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `lifted_at` (String) When someone lifted the block.
- `project_id` (String) ID of the project this block belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
