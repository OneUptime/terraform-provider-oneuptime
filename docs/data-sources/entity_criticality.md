---
page_title: "oneuptime_entity_criticality Data Source - oneuptime"
subcategory: "Security"
description: |-
  How much each user, host, IP address and data store matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).
---

# oneuptime_entity_criticality (Data Source)

How much each user, host, IP address and data store matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one entity criticality may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_entity_criticality" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_entity_criticality" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the entry. The ID of a `oneuptime_user` (see the data source).
- `criticality` (String) How much it matters: Low (its risk counts half), High (one and a half times) or Extreme (twice).
- `description` (String) Why it matters this much, for whoever reads the list later.
- `entity_type` (String) User, Host, IP Address or Data Store: the kind of entity it is for. Empty is any kind.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `value` (String) The user's sign-in name, the host's name, the address or the data store, matched whatever the case. * stands for any run of characters: dc-* names every domain controller.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
