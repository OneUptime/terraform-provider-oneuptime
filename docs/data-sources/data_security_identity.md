---
page_title: "oneuptime_data_security_identity Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Every identity that can reach the project's data stores, with what it reaches, how sensitive that is, and how much harm it could do.
---

# oneuptime_data_security_identity (Data Source)

Every identity that can reach the project's data stores, with what it reaches, how sensitive that is, and how much harm it could do.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security identity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_identity" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_security_identity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `can_change_sensitive_data` (Boolean) It may write to, or control, a store holding Confidential or Restricted data.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `estimated_records` (Number) About how many records with sensitive data the stores it reaches hold, where scans could tell.
- `highest_access_level` (String) The most it may do to any store it reaches: Read, Write or Admin.
- `highest_sensitivity` (String) How sensitive the most sensitive store it reaches is: Restricted, Confidential, Internal or Public.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `identity_key` (String) The principal as its cloud writes it, which identifies the identity within its project.
- `is_external` (Boolean) The identity is outside the project's connected accounts.
- `is_public` (Boolean) The identity is anyone on the internet, or any AWS account.
- `is_seen_in_traces` (Boolean) OneUptime saw the identity use a store in traces.
- `is_unused` (Boolean) The identity has not been used in 90 days, as far as its cloud can say.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `principal` (String) Who or what it is, as its cloud names it: *, an account id, an ARN, an email, a service account.
- `principal_name` (String) The identity's name, for people.
- `principal_type` (String) What it is: Public, Account, Role, User, Group, Service Principal, Service Account or Service.
- `provider_value` (String) Where the identity lives: AWS, Azure, Google Cloud, or OneUptime for a service seen in traces.
- `restricted_store_count` (Number) How many of them hold Restricted data.
- `risk_score` (Number) How much harm the identity could do with what it reaches, from 0 to 100.
- `sensitive_store_count` (Number) How many of them hold Confidential or Restricted data.
- `store_count` (Number) How many data stores the identity can reach.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `data_type_ids` (String) The data types found in the sensitive stores it reaches, by their ids. A JSON value: read it with `jsondecode()`.
- `first_seen_at` (String) When the identity was first found reaching a store.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_seen_at` (String) When the identity was last found reaching a store.
- `last_used_at` (String) When the identity was last used, where OneUptime can tell.
- `project_id` (String) ID of the project this identity belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
