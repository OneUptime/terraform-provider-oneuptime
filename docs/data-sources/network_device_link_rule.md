---
page_title: "oneuptime_network_device_link_rule Data Source - oneuptime"
subcategory: "Network"
description: |-
  Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.
---

# oneuptime_network_device_link_rule (Data Source)

Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device link rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_link_rule" "example" {
  name = "Example network device link rule"
}

# Or by id:
data "oneuptime_network_device_link_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule draws links. Disable to take its edges off the map without deleting the rule.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this rule.
- `scope` (String) How wide the 'exactly one parent device' question is asked. Project (the default) looks for one parent across the whole project. Site asks once per site, so the same rule can draw an uplink in every building. Rules created before this existed are Project.

### Read-Only

- `child_device_labels` (Set of String) Devices carrying ALL of these labels each get one uplink drawn to the parent device. Empty matches nothing — a rule that linked every device in the project is never what anyone meant. IDs of `oneuptime_label` resources.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `parent_device_labels` (Set of String) The device carrying ALL of these labels is what the children uplink to. It has to identify exactly one device: match none and the rule draws nothing, match several and the rule is ambiguous and also draws nothing. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
