---
page_title: "oneuptime_runbook Resource - oneuptime"
subcategory: "Other"
description: |-
  Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.
---

# oneuptime_runbook (Resource)

Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.

## Example Usage

```terraform
resource "oneuptime_runbook" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Is this runbook enabled?..
- `steps` (String) Ordered list of steps to run for this runbook. Each step is one of Manual, JavaScript, HTTP request, Bash or AI...
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook.example <id>
```
