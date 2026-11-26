---
page_title: "oneuptime_on_call_schedule_layer Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layers
---

# oneuptime_on_call_schedule_layer (Resource)

On-Call Schedule Layers

## Example Usage

```terraform
# Weekdays, 09:00 to 17:00. Times are wall-clock times in the schedule's
# timezone; these are for a schedule in UTC.
resource "oneuptime_on_call_schedule_layer" "example" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
  name                            = "Business hours"
  description                     = "Weekdays, 09:00 to 17:00"
  order                           = 1 # Lower wins: this layer is on call in its hours.
  starts_at                       = "2026-01-05T09:00:00Z"
  hand_off_time                   = "2026-01-12T09:00:00Z"

  # The next person takes over every week.
  rotation = jsonencode({
    _type = "Recurring"
    value = {
      intervalType  = "Week"
      intervalCount = { _type = "PositiveNumber", value = 1 }
    }
  })

  # One window a day, each dated on its weekday (2026-01-05 is a Monday).
  restriction_times = jsonencode({
    _type = "RestrictionTimes"
    value = {
      restictionType = "Weekly"
      weeklyRestrictionTimes = [
        for day, date in {
          Monday    = "2026-01-05"
          Tuesday   = "2026-01-06"
          Wednesday = "2026-01-07"
          Thursday  = "2026-01-08"
          Friday    = "2026-01-09"
          } : {
          startDay  = day
          endDay    = day
          startTime = { _type = "DateTime", value = "${date}T09:00:00Z" }
          endTime   = { _type = "DateTime", value = "${date}T17:00:00Z" }
        }
      ]
    }
  })
}

# Nights and weekends: on call around the clock, below the business-hours
# layer, so it covers every hour that layer does not.
resource "oneuptime_on_call_schedule_layer" "around_the_clock" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
  name                            = "Nights and weekends"
  order                           = 2
  starts_at                       = "2026-01-05T09:00:00Z"
  hand_off_time                   = "2026-01-06T09:00:00Z"

  rotation = jsonencode({
    _type = "Recurring"
    value = {
      intervalType  = "Day"
      intervalCount = { _type = "PositiveNumber", value = 1 }
    }
  })

  restriction_times = jsonencode({
    _type = "RestrictionTimes"
    value = { restictionType = "None" }
  })
}
```

<a id="rotation-and-restriction-times"></a>
## Rotation and restriction times

`rotation` and `restriction_times` take JSON: write them with `jsonencode()`, in the shapes below. The provider checks every `_type` and timestamp in them at plan time.

Both are optional. Left out, a layer hands off every week and is on call around the clock (OneUptime 14.0.73 and later; earlier versions refuse a layer without them).

### Rotation

How long each turn lasts. The first hand-off is at `hand_off_time`, and the duty moves on to the next person every rotation after that.

```terraform
rotation = jsonencode({
  _type = "Recurring"
  value = {
    intervalType  = "Week" # Hour, Day, Week, Month or Year
    intervalCount = { _type = "PositiveNumber", value = 1 }
  }
})
```

`intervalCount` is a whole number, 1 or more: with `value = 2` and `intervalType = "Week"`, the duty hands off every two weeks.

### Restriction times

When the layer is on call. `restictionType` is the key's real name, without the second "r": the server ignores a key spelled any other way.

- `None`: around the clock.

  ```terraform
  restriction_times = jsonencode({
    _type = "RestrictionTimes"
    value = { restictionType = "None" }
  })
  ```

- `Daily`: the same hours every day. Only the time of day of `startTime` and `endTime` counts, and a window that ends before it starts runs past midnight.

  ```terraform
  restriction_times = jsonencode({
    _type = "RestrictionTimes"
    value = {
      restictionType = "Daily"
      dayRestrictionTimes = {
        startTime = { _type = "DateTime", value = "2026-01-05T09:00:00Z" }
        endTime   = { _type = "DateTime", value = "2026-01-05T17:00:00Z" }
      }
    }
  })
  ```

- `Weekly`: windows in the week. A window runs from its `startTime` to its `endTime`, and both the weekday and the time of day count: date each one on the weekday you mean, and name those weekdays in `startDay` and `endDay`. A window is continuous, so Monday 09:00 to Friday 17:00 includes the nights in between. For working hours, give each day a window of its own, as the example above does. A window that ends earlier in the week than it starts runs over the weekend:

  ```terraform
  restriction_times = jsonencode({
    _type = "RestrictionTimes"
    value = {
      restictionType = "Weekly"
      weeklyRestrictionTimes = [{
        startDay  = "Friday"
        endDay    = "Monday"
        startTime = { _type = "DateTime", value = "2026-01-09T17:00:00Z" } # a Friday
        endTime   = { _type = "DateTime", value = "2026-01-12T09:00:00Z" } # the Monday after
      }]
    }
  })
  ```

A time is `{ _type = "DateTime", value = "<RFC 3339 timestamp>" }`, with or without fractional seconds. It is a wall-clock time in the schedule's `timezone`: write it with the offset that time zone has on that date, as in `2026-01-05T09:00:00-05:00` for 09:00 in New York. A schedule without a `timezone` reads times in the server's time zone.

<a id="order"></a>
## Order

`order` is this layer's place among the layers of its schedule, the lowest number the highest priority. Setting it places this layer and moves no other, so a configuration that numbers every layer of a schedule is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.layers`:

```terraform
order = index(local.layers, each.key) + 1
```

- Give each layer a number of its own: two layers at one number are ranked in the order of their IDs.
- A layer created without `order` goes below the schedule's last layer.
- Deleting a layer leaves a gap in the numbers, which changes nothing: layers go by their numbers, not by counting.
- Moving a layer in the OneUptime Dashboard numbers the schedule's layers 1..n again; the next plan shows that as changes to `order`.

To set the layers' order and who takes turns in each in one write, use [`oneuptime_on_call_schedule_layers`](./on_call_schedule_layers) instead, and leave `order` unset here.


## Schema

### Required

- `hand_off_time` (String) Hand off time. When would you like to hand off the duty to the next user in this layer?
- `name` (String) Friendly name for this layer.
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the layer is in. The ID of a `oneuptime_on_call_policy_schedule`.
- `starts_at` (String) Start date and time of this layer.

### Optional

- `description` (String) Description for this layer. This is optional and can be left blank.
- `order` (Number) This layer's priority in its schedule: the lower the number, the higher the priority. Layers that share a number are ranked in the order of their IDs. Setting it places this layer and moves no other, so a configuration that numbers every layer of a schedule ends up exactly as written. Left out when you create a layer, it goes below the schedule's last layer. See [Order](#order).
- `restriction_times` (String) Restrict this layer to these times. A JSON value: write it with `jsonencode()`. See [Rotation and restriction times](#rotation-and-restriction-times).
- `rotation` (String) How often would you like to hand off the duty to the next user in this layer? A JSON value: write it with `jsonencode()`. See [Rotation and restriction times](#rotation-and-restriction-times).

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

Import an existing on call schedule layer by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_layer.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layer.example <id>
```
