---
page_title: "oneuptime_entity_criticality Data Source - oneuptime"
subcategory: "Other"
description: |-
  How much each user, host and IP address matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).
---

# oneuptime_entity_criticality (Data Source)

How much each user, host and IP address matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one entity criticality may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_entity_criticality" "example" {
  value = "example-value"
}

# Or by id:
data "oneuptime_entity_criticality" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who added the entry. The ID of a `oneuptime_user` (see the data source).
- `criticality` (String) How much it matters: Low (its risk counts half), High (one and a half times) or Extreme (twice).
- `description` (String) Why it matters this much, for whoever reads the list later.
- `entity_type` (String) User, Host or IP Address: the kind of entity it is for. Empty is any kind.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `value` (String) The user's sign-in name, the host's name or the address, matched whatever the case. * stands for any run of characters: dc-* names every domain controller.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
