---
page_title: "oneuptime_block Resource - oneuptime"
subcategory: "Other"
description: |-
  The IP addresses and ranges the project blocks: why, until when, and where each block is enforced.
---

# oneuptime_block (Resource)

The IP addresses and ranges the project blocks: why, until when, and where each block is enforced.

## Example Usage

```terraform
resource "oneuptime_block" "example" {
  value = "35.205.11.131"
}
```

## Schema

### Required

- `value` (String) The address or range blocked: an IPv4 or IPv6 address, or a range such as 203.0.113.0/24. Kept as its first address and its length, an IPv6 address compressed in lowercase, a /32 or /128 as the bare address.

### Optional

- `action` (String) Block turns the address's requests away. Challenge asks for a challenge a person passes and a script or an API client does not. Allowed values: `Block`, `Challenge`. Defaults to `Block`.
- `expires_at` (String) When the block ends by itself: 24 hours after it is made unless the request says otherwise. Null keeps it until someone lifts it, which takes a Security Admin.
- `lift_reason` (String) Why the block was lifted, given with the lift.
- `reason` (String) Why the address is blocked, for whoever reads the list later.
- `security_case_id` (String) ID of the case the block was made for. The block and its lifting go on that case's timeline. The ID of a `oneuptime_security_case`.
- `state` (String) Active until it expires or someone lifts it. Set it to Lifted to lift the block; Expired is set by OneUptime when its time is up. Allowed values: `Active`, `Lifted`, `Expired`. Defaults to `Active`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who blocked it. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `indicator_type` (String) What the block holds: IpAddress for one address, IpRange for a range. Worked out from the value.
- `lifted_at` (String) When someone lifted the block.
- `lifted_by_user_id` (String) ID of the person who lifted the block. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this block belongs to. The ID of a `oneuptime_project`.
- `security_playbook_run_id` (String) ID of the playbook run whose step made the block. The ID of a `oneuptime_security_playbook_run`.
- `source` (String) Who or what made the block: a Person, a Playbook, an AutomaticResponse policy, OneUptimeAi or an Api key.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing block by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_block.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_block.example <id>
```
