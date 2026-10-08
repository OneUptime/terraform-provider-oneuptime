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
  name        = "Example on call policy"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived on-call policies are hidden from the On-Call Policies list and page no one: incidents and alerts that use them skip them. Unarchiving puts them back in service. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `repeat_policy_if_no_one_acknowledges` (Boolean) Repeat the policy if no one acknowledges the alert. Defaults to `false`.
- `repeat_policy_if_no_one_acknowledges_no_of_times` (Number) Repeat the policy X number of times if no one acknowledges the alert. Defaults to `0`.

### Read-Only

- `archived_at` (String) When this on-call policy was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call policy by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_policy.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_policy.example <id>
```
