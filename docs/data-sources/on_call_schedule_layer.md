---
page_title: "oneuptime_on_call_schedule_layer Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layers
---

# oneuptime_on_call_schedule_layer (Data Source)

On-Call Schedule Layers

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call schedule layer may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_schedule_layer" "example" {
  name = "Example on call schedule layer"
}

# Or by id:
data "oneuptime_on_call_schedule_layer" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description for this layer. This is optional and can be left blank.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this layer.
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the layer is in. The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) This layer's priority in its schedule: the lower the number, the higher the priority. Layers that share a number are ranked in the order of their IDs. Setting it places this layer and moves no other, so a configuration that numbers every layer of a schedule ends up exactly as written. Left out when you create a layer, it goes below the schedule's last layer.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `hand_off_time` (String) Hand off time. When would you like to hand off the duty to the next user in this layer?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `restriction_times` (String) Restrict this layer to these times. A JSON value: read it with `jsondecode()`.
- `rotation` (String) How often would you like to hand off the duty to the next user in this layer? A JSON value: read it with `jsondecode()`.
- `starts_at` (String) Start date and time of this layer.
- `updated_at` (String) Date and Time when the object was updated.
