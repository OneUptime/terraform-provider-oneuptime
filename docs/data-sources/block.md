---
page_title: "oneuptime_block Data Source - oneuptime"
subcategory: "Other"
description: |-
  The IP addresses and ranges the project blocks: why, until when, and where each block is enforced.
---

# oneuptime_block (Data Source)

The IP addresses and ranges the project blocks: why, until when, and where each block is enforced.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one block may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_block" "example" {
  indicator_type = "IpAddress"
}

# Or by id:
data "oneuptime_block" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) Block turns the address's requests away. Challenge asks for a challenge a person passes and a script or an API client does not.
- `created_by_user_id` (String) ID of the person who blocked it. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `indicator_type` (String) What the block holds: IpAddress for one address, IpRange for a range. Worked out from the value.
- `lift_reason` (String) Why the block was lifted, given with the lift.
- `lifted_by_user_id` (String) ID of the person who lifted the block. The ID of a `oneuptime_user` (see the data source).
- `reason` (String) Why the address is blocked, for whoever reads the list later.
- `security_case_id` (String) ID of the case the block was made for. The block and its lifting go on that case's timeline. The ID of a `oneuptime_security_case`.
- `security_playbook_run_id` (String) ID of the playbook run whose step made the block. The ID of a `oneuptime_security_playbook_run`.
- `source` (String) Who or what made the block: a Person, a Playbook, an AutomaticResponse policy, OneUptimeAi or an Api key.
- `state` (String) Active until it expires or someone lifts it. Set it to Lifted to lift the block; Expired is set by OneUptime when its time is up.
- `value` (String) The address or range blocked: an IPv4 or IPv6 address, or a range such as 203.0.113.0/24. Kept as its first address and its length, an IPv6 address compressed in lowercase, a /32 or /128 as the bare address.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `expires_at` (String) When the block ends by itself: 24 hours after it is made unless the request says otherwise. Null keeps it until someone lifts it, which takes a Security Admin.
- `lifted_at` (String) When someone lifted the block.
- `project_id` (String) ID of the project this block belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
