---
page_title: "oneuptime_runbook Data Source - oneuptime"
subcategory: "Other"
description: |-
  Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.
---

# oneuptime_runbook (Data Source)

Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runbook may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runbook" "example" {
  name = "Example runbook"
}

# Or by id:
data "oneuptime_runbook" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Is this runbook enabled?
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `steps` (String) Ordered list of steps to run for this runbook. Each step is one of Manual, JavaScript, HTTP request, Bash or AI. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
