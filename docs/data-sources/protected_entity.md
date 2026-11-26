---
page_title: "oneuptime_protected_entity Data Source - oneuptime"
subcategory: "Security"
description: |-
  Users, hosts and IP addresses that security playbooks never contain and the blocklist never blocks, such as break-glass accounts, domain controllers, VPN gateways and office addresses.
---

# oneuptime_protected_entity (Data Source)

Users, hosts and IP addresses that security playbooks never contain and the blocklist never blocks, such as break-glass accounts, domain controllers, VPN gateways and office addresses.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one protected entity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_protected_entity" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_protected_entity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who protected it. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Why it is protected, for whoever reads the list later.
- `entity_type` (String) User, Host or IP Address: what the value names.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `value` (String) The user's sign-in name or the host's name, matched whatever the case. * stands for any run of characters: svc-* protects every service account. For an IP Address, an IPv4 or IPv6 address or a range such as 203.0.113.0/24, which protects every address inside it.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this protected entity belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
