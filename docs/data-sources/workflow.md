---
page_title: "oneuptime_workflow Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Integrate your OneUptime project with rest of your software stack.
---

# oneuptime_workflow (Data Source)

Integrate your OneUptime project with rest of your software stack. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workflow" "by_name" {
  name = "example-workflow"
}

data "oneuptime_workflow" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Archived workflows are hidden from the Workflows list and never run, from any trigger. Unarchiving restores them as they were... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Is this workflow enabled?.. Computed.
- `graph` (String) Workflow Graph in JSON. Ideally, create this via UI and not via API... Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `webhook_secret_key` (String) The secret part of the Webhook trigger's URL (/workflow/trigger/<key>). Anyone who has the URL can start the workflow, so only people who can edit the workflow can read the key. Generated when the workflow is created; set a new value to reset the URL... Computed.
- `incoming_email_secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
