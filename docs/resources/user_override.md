---
page_title: "oneuptime_user_override Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.
---

# oneuptime_user_override (Resource)

While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.

## Example Usage

```terraform
resource "oneuptime_user_override" "example" {
  override_user_id = "123e4567-e89b-12d3-a456-426614174000"
  route_alerts_to_user_id = "123e4567-e89b-12d3-a456-426614174000"
  starts_at = "2030-01-01T00:00:00Z"
  ends_at = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `override_user_id` (String) A unique identifier for an object, represented as a UUID..
- `route_alerts_to_user_id` (String) A unique identifier for an object, represented as a UUID..
- `starts_at` (String) A date time object..
- `ends_at` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_user_override.example <id>
```
