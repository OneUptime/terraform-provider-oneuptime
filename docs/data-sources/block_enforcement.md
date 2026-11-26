---
page_title: "oneuptime_block_enforcement Data Source - oneuptime"
subcategory: "Other"
description: |-
  Where each block of the blocklist is enforced, and whether that holds: one row per block and destination, written by OneUptime as it applies the block and reads it back.
---

# oneuptime_block_enforcement (Data Source)

Where each block of the blocklist is enforced, and whether that holds: one row per block and destination, written by OneUptime as it applies the block and reads it back.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one block enforcement may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_block_enforcement" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_block_enforcement" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `connection_id` (String) ID of the connection the destination is enforced through: a security connection, a cloud connection or a feed.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `destination_key` (String) Which destination this is, as OneUptime keys it: its provider and its connection, and what in it (a zone, a web ACL) when it holds more than one.
- `destination_name` (String) The destination as people read it: Cloudflare, AWS WAF (prod-alb).
- `error` (String) Why the block is not applied there: the destination cannot challenge, matches IPv4 only, is full, or the provider refused.
- `failed_attempts` (Number) How many tries in a row failed at the destination; the next waits longer each time.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `provider_ref` (String) The provider's own ID of what holds the block: a list item, a rule, an indicator.
- `provider_value` (String) The kind of destination: Cloudflare, AWS WAF, Azure WAF, Cloud Armor, a feed, Microsoft Defender or CrowdStrike.
- `security_block_id` (String) ID of the block this row says the enforcement of. The ID of a `oneuptime_block`.
- `state` (String) Pending until the destination is seen holding the block, Applied once it is, Failed with why it was not applied, then Removing and Removed when the block ends.

### Read-Only

- `applied_at` (String) When the destination was first seen holding the block.
- `created_at` (String) Date and Time when the object was created.
- `error_reason` (String) Why the block is not applied there, as a sentence and its values that the Dashboard words in the reader's language. A JSON value: read it with `jsondecode()`.
- `last_confirmed_at` (String) When OneUptime last read the destination back and found the block there.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this row belongs to. The ID of a `oneuptime_project`.
- `retry_at` (String) When OneUptime tries the destination again after a failure; empty when nothing waits.
- `updated_at` (String) Date and Time when the object was updated.
