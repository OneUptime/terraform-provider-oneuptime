---
page_title: "oneuptime_data_store_access_entry Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Who and what can reach each store: the internet, other accounts, identities and the services seen using it, with the policy that grants it.
---

# oneuptime_data_store_access_entry (Data Source)

Who and what can reach each store: the internet, other accounts, identities and the services seen using it, with the policy that grants it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store access entry may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store_access_entry" "example" {
  created_by_actor_type = "User"
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
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_store_id` (String) The store that can be reached. The ID of a `oneuptime_data_store`.
- `granted_via` (String) What grants the access: a bucket policy, an ACL, a role assignment, an IAM binding, a snapshot share, or traffic seen in traces.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_external` (Boolean) The principal is outside the project's connected accounts.
- `is_public` (Boolean) The principal is anyone on the internet.
- `is_seen_in_traces` (Boolean) OneUptime saw this principal use the store in traces.
- `is_unused` (Boolean) The principal has not been used in 90 days, as far as its cloud can say.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `principal` (String) Who or what has access, as its cloud names it: *, an account id, an ARN, an email, a service account.
- `principal_name` (String) The principal's name, for people.
- `principal_type` (String) What the principal is: Public, Account, Role, User, Group, Service Principal, Service Account or Service.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When the access was first read.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_seen_at` (String) When the access was last read.
- `last_used_at` (String) When the access was last used, where OneUptime can tell.
- `project_id` (String) ID of the project this access entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
