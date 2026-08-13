---
page_title: "oneuptime_on_call_policy Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty, schedules and roster for your project
---

# oneuptime_on_call_policy (Resource)

Manage on-call duty, schedules and roster for your project

## Example Usage

```terraform
resource "oneuptime_on_call_policy" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `repeat_policy_if_no_one_acknowledges` (Bool) Repeat the policy if no one acknowledges the alert..
- `repeat_policy_if_no_one_acknowledges_no_of_times` (Number) Repeat the policy X number of times if no one acknowledges the alert..
- `custom_fields` (String) Custom Fields on this resource...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_on_call_policy.example <id>
```
