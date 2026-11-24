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
  security_block_id = oneuptime_block.example.id
}

# Or by id:
data "oneuptime_block_enforcement" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `connection_id` (String) ID of the connection the destination is enforced through: a security connection, a cloud connection or a feed.
- `destination_key` (String) Which destination this is, as OneUptime keys it: its provider and its connection, and what in it (a zone, a web ACL) when it holds more than one.
- `destination_name` (String) The destination as people read it: Cloudflare, AWS WAF (prod-alb).
- `error` (String) Why the block is not applied there: the destination cannot challenge, matches IPv4 only, is full, or the provider refused.
- `failed_attempts` (Number) How many tries in a row failed at the destination; the next waits longer each time.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `provider_ref` (String) The provider's own ID of what holds the block: a list item, a rule, an indicator.
- `provider_value` (String) The kind of destination: Cloudflare, AWS WAF, Azure WAF, Cloud Armor, a feed, Microsoft Defender or CrowdStrike.
- `security_block_id` (String) ID of the block this row says the enforcement of. The ID of a `oneuptime_block`.
- `state` (String) Pending until the destination is seen holding the block, Applied once it is, Failed with why it was not applied, then Removing and Removed when the block ends.

### Read-Only

- `applied_at` (String) When the destination was first seen holding the block.
- `created_at` (String) Date and Time when the object was created.
- `last_confirmed_at` (String) When OneUptime last read the destination back and found the block there.
- `project_id` (String) ID of the project this row belongs to. The ID of a `oneuptime_project`.
- `retry_at` (String) When OneUptime tries the destination again after a failure; empty when nothing waits.
- `updated_at` (String) Date and Time when the object was updated.
