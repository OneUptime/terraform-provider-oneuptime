---
page_title: "oneuptime_workflow Resource - oneuptime"
subcategory: "Workflows"
description: |-
  Integrate your OneUptime project with rest of your software stack.
---

# oneuptime_workflow (Resource)

Integrate your OneUptime project with rest of your software stack.

## Example Usage

```terraform
resource "oneuptime_workflow" "example" {
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
- `is_enabled` (Bool) Is this workflow enabled?..
- `graph` (String) Workflow Graph in JSON. Ideally, create this via UI and not via API...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `webhook_secret_key` (String) Secret key used to trigger this workflow via webhook. Use this instead of the workflow ID for security...

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
terraform import oneuptime_workflow.example <id>
```
