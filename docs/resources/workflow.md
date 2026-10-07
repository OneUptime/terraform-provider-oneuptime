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
- `is_archived` (Bool) Archived workflows are hidden from the Workflows list and never run, from any trigger. Unarchiving restores them as they were...
- `is_enabled` (Bool) Is this workflow enabled?..
- `graph` (String) Workflow Graph in JSON. Ideally, create this via UI and not via API...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `webhook_secret_key` (String) The secret part of the Webhook trigger's URL (/workflow/trigger/<key>). Anyone who has the URL can start the workflow, so only people who can edit the workflow can read the key. Generated when the workflow is created; set a new value to reset the URL...
- `incoming_email_secret_key` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_workflow.example <id>
```
