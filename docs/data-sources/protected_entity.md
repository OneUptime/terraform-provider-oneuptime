---
page_title: "oneuptime_protected_entity Data Source - oneuptime"
subcategory: "Other"
description: |-
  Users and hosts that security playbooks never contain, such as break-glass accounts, domain controllers and VPN gateways.
---

# oneuptime_protected_entity (Data Source)

Users and hosts that security playbooks never contain, such as break-glass accounts, domain controllers and VPN gateways.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one protected entity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_protected_entity" "example" {
  entity_type = "example-entity-type"
}

# Or by id:
data "oneuptime_protected_entity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who protected it. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Why it is protected, for whoever reads the list later.
- `entity_type` (String) User or Host: what the value names.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `value` (String) The user's sign-in name or the host's name, matched whatever the case. * stands for any run of characters: svc-* protects every service account.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this protected entity belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
