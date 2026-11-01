---
page_title: "oneuptime_data_store_access_entry Data Source - oneuptime"
subcategory: "Other"
description: |-
  Who and what can reach each store: the internet, other accounts, identities and the services seen using it, with the policy that grants it.
---

# oneuptime_data_store_access_entry (Data Source)

Who and what can reach each store: the internet, other accounts, identities and the services seen using it, with the policy that grants it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store access entry may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store_access_entry" "example" {
  data_store_id = oneuptime_data_store.example.id
}

# Or by id:
data "oneuptime_data_store_access_entry" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `access_key` (String) The principal, level and grant, which identify the entry within its store.
- `access_level` (String) What it may do: Read, Write or Admin.
- `data_store_id` (String) The store that can be reached. The ID of a `oneuptime_data_store`.
- `granted_via` (String) What grants the access: a bucket policy, an ACL, a role assignment, an IAM binding, a snapshot share, or traffic seen in traces.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_external` (Boolean) The principal is outside the project's connected accounts.
- `is_public` (Boolean) The principal is anyone on the internet.
- `is_seen_in_traces` (Boolean) OneUptime saw this principal use the store in traces.
- `principal` (String) Who or what has access, as its cloud names it: *, an account id, an ARN, an email, a service account.
- `principal_name` (String) The principal's name, for people.
- `principal_type` (String) What the principal is: Public, Account, Role, User, Group, Service Principal, Service Account or Service.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When the access was first read.
- `last_seen_at` (String) When the access was last read.
- `last_used_at` (String) When the access was last used, where OneUptime can tell.
- `project_id` (String) ID of the project this access entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
