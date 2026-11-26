---
page_title: "oneuptime_network_site_assignment_rule Data Source - oneuptime"
subcategory: "Network"
description: |-
  Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.
---

# oneuptime_network_site_assignment_rule (Data Source)

Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site assignment rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_assignment_rule" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_network_site_assignment_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `hostname_pattern` (String) Devices whose hostname, SNMP system name, display name or DNS name matches this wildcard pattern are assigned to the site.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `priority` (Number) Where this rule sits in the list: when several rules match a device, the one highest in the list wins. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number places this rule there and moves no other rule: two rules at one number are read the older first. In the dashboard, drag the rows to reorder them.
- `site_id` (String) ID of the Network Site matched resources are assigned to. The ID of a `oneuptime_network_site`.
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
