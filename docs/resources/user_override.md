---
page_title: "oneuptime_user_override Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. With nobody named to cover, it skips them instead: each of their turns in an on-call schedule goes to the next person in the rotation. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.
---

# oneuptime_user_override (Resource)

While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. With nobody named to cover, it skips them instead: each of their turns in an on-call schedule goes to the next person in the rotation. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.

## Example Usage

```terraform
resource "oneuptime_user_override" "example" {
  override_user_id = data.oneuptime_user.example.id
  starts_at        = "2030-01-01T00:00:00Z"
  ends_at          = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `ends_at` (String) When the override ends, which has to be after it starts. From then on, alerts page the Override User again, and their turns in each rotation are theirs again.
- `override_user_id` (String) ID of the user who is away. While the override is in force, alerts that would page this user go to the user in Route Alerts To User ID instead. With no Route Alerts To User ID, they are skipped: each of their turns in an on-call schedule goes to the next person in the rotation who is not away. The ID of a `oneuptime_user` (see the data source).
- `starts_at` (String) When the override starts sending the Override User's alerts to the Route Alerts To User - or, with none, skipping their turns.

### Optional

- `on_call_duty_policy_id` (String) ID of the on-call policy this override applies to. Leave it empty for a global override, which applies to every on-call policy. The ID of a `oneuptime_on_call_policy`.
- `route_alerts_to_user_id` (String) ID of the user who covers. While the override is in force, this user gets the alerts that would page the user in Override User ID. Leave it empty to skip that user instead: each of their turns in an on-call schedule goes to the next person in the rotation who is not away, and teams and escalation rules that name them do not page them until the override ends. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing user override by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_user_override.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_user_override.example <id>
```
