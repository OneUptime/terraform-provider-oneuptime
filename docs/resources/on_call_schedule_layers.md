---
page_title: "oneuptime_on_call_schedule_layers Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Who takes turns in each layer of an on-call schedule, in turn order, and the order of the layers, saved all at once.
---

# oneuptime_on_call_schedule_layers (Resource)

Who takes turns in each layer of an on-call schedule, in turn order, and the order of the layers, saved all at once.

## Example Usage

```terraform
resource "oneuptime_on_call_policy_schedule" "sre" {
  name     = "SRE"
  timezone = "UTC"
}

resource "oneuptime_on_call_schedule_layer" "primary" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.sre.id
  name                            = "Primary"
  starts_at                       = "2030-01-07T09:00:00Z"
  hand_off_time                   = "2030-01-14T09:00:00Z"
}

resource "oneuptime_on_call_schedule_layer" "fallback" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.sre.id
  name                            = "Fallback"
  starts_at                       = "2030-01-07T09:00:00Z"
  hand_off_time                   = "2030-01-14T09:00:00Z"
}

data "oneuptime_user" "alice" {
  email = "alice@example.com"
}

data "oneuptime_user" "bob" {
  email = "bob@example.com"
}

data "oneuptime_user" "carol" {
  email = "carol@example.com"
}

# Alice and Bob take turns on top; Carol covers whenever neither is on call.
resource "oneuptime_on_call_schedule_layers" "sre" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.sre.id

  layers = [
    {
      layer_id = oneuptime_on_call_schedule_layer.primary.id
      user_ids = [data.oneuptime_user.alice.id, data.oneuptime_user.bob.id]
    },
    {
      layer_id = oneuptime_on_call_schedule_layer.fallback.id
      user_ids = [data.oneuptime_user.carol.id]
    },
  ]
}
```

## How a change is saved

Each apply sends the whole change in one call, as **Save changes** on a schedule's **Layers** page in OneUptime does: OneUptime writes all of it or none of it, and works out who is on call once, at once, so someone who goes on or off call, or is next, hears it once. Someone who stays in a layer keeps their place in it, and only their turn moves.

The change is made to the layers as they are read just before it. If someone changes them in between, OneUptime refuses the save, and the provider reads them again and saves again.

## Layers it does not list

List layers of the schedule by their `oneuptime_on_call_schedule_layer` IDs, the highest priority first. A layer of the schedule that is not listed keeps its people and its place, and the layers listed take the places they have between them, in the order listed. To delete a layer, remove its `oneuptime_on_call_schedule_layer` and take it out of `layers` in the same change.

Give the people of the layers it lists here only, not with `oneuptime_on_call_schedule_layer_user` too, and leave those layers' `order` unset: this resource orders them.

## Destroying it

Destroying it takes everyone out of the layers it lists, and leaves the layers. To stop managing their people with Terraform and keep them, remove it from the state instead, with a `removed` block or `terraform state rm`.

## Schema

### Required

- `layers` (Attributes List) Layers of the schedule, the highest priority first: the layer on top is the one that pages while someone in it is on call. Each layer once. The layers listed take the places they have in the schedule between them, in this order. (see [below for nested schema](#nestedatt--layers))
- `on_call_duty_policy_schedule_id` (String) ID of the `oneuptime_on_call_policy_schedule` whose layers these are. Changing it replaces the resource.

### Read-Only

- `id` (String) The schedule's ID: a schedule has one of these.

<a id="nestedatt--layers"></a>
### Nested Schema for `layers`

Required:

- `layer_id` (String) ID of the `oneuptime_on_call_schedule_layer`: a layer of this schedule.
- `user_ids` (List of String) User IDs of the people who take turns in the layer, in turn order: `[]` for nobody. Someone in the project, or invited to it - `oneuptime_team_member`'s `user_id`, or `data.oneuptime_user`'s `id` for someone named by email address.

## Import

Import the layers of an existing schedule by the schedule's ID, with an `import` block (Terraform 1.5+, OpenTofu 1.6+). The import takes every layer of the schedule, in its order:

```terraform
import {
  to = oneuptime_on_call_schedule_layers.example
  id = "<on_call_duty_policy_schedule_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layers.example <on_call_duty_policy_schedule_id>
```
