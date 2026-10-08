---
page_title: "oneuptime_on_call_policy Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty, schedules and roster for your project
---

# oneuptime_on_call_policy (Data Source)

Manage on-call duty, schedules and roster for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call policy may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_policy" "example" {
  name = "Example on call policy"
}

# Or by id:
data "oneuptime_on_call_policy" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Archived on-call policies are hidden from the On-Call Policies list and page no one: incidents and alerts that use them skip them. Unarchiving puts them back in service.
- `name` (String) Any friendly name of this object.
- `repeat_policy_if_no_one_acknowledges` (Boolean) Repeat the policy if no one acknowledges the alert.
- `repeat_policy_if_no_one_acknowledges_no_of_times` (Number) Repeat the policy X number of times if no one acknowledges the alert.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When this on-call policy was archived. Empty while it is not archived.
- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
