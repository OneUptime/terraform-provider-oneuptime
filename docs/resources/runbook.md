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
  name        = "Example runbook"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_enabled` (Boolean) Is this runbook enabled? Defaults to `true`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `steps` (String) Ordered list of steps to run for this runbook. Each step is one of Manual, JavaScript, HTTP request, Bash or AI. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook.example <id>
```
