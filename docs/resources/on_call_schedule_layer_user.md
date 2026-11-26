---
page_title: "oneuptime_on_call_schedule_layer_user Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users to the layers of your on-call schedules.
---

# oneuptime_on_call_schedule_layer_user (Resource)

Add users to the layers of your on-call schedules.

## Example Usage

```terraform
resource "oneuptime_on_call_schedule_layer_user" "example" {
  on_call_duty_policy_schedule_id       = oneuptime_on_call_policy_schedule.example.id
  on_call_duty_policy_schedule_layer_id = oneuptime_on_call_schedule_layer.example.id
  user_email                            = "alice@example.com"
}
```

<a id="order"></a>
## Order

`order` is this person's place among the people who take turns in its layer, the lowest first. Setting it places this person and moves no other, so a configuration that numbers every person of a layer is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.people`:

```terraform
order = index(local.people, each.key) + 1
```

- Give each person a number of its own: two people at one number take their turns in the order of their IDs.
- A person created without `order` takes the layer's last turn.
- Deleting a person leaves a gap in the numbers, which changes nothing: people go by their numbers, not by counting.
- Moving a person in the OneUptime Dashboard numbers the layer's people 1..n again; the next plan shows that as changes to `order`.

To set the layers' order and who takes turns in each in one write, use [`oneuptime_on_call_schedule_layers`](./on_call_schedule_layers) instead, and leave `order` unset here.


## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the layer is in. The ID of a `oneuptime_on_call_policy_schedule`.
- `on_call_duty_policy_schedule_layer_id` (String) ID of the layer the person takes turns on call in. The ID of a `oneuptime_on_call_schedule_layer`.

### Optional

- `order` (Number) This person's turn in the layer's rotation, lowest first. People who share a turn take it in the order of their IDs. Setting it places this person and moves no other, so a configuration that numbers everyone in a layer ends up exactly as written. Left out when you add someone, they take the layer's last turn. See [Order](#order).
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person who takes turns on call in this layer. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call schedule layer user by its on call schedule layer and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_layer_user.example
  id = "<on_call_duty_policy_schedule_layer_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layer_user.example <on_call_duty_policy_schedule_layer_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
