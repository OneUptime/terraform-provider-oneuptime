---
page_title: "oneuptime_workflow Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Integrate your OneUptime project with rest of your software stack.
---

# oneuptime_workflow (Data Source)

Integrate your OneUptime project with rest of your software stack.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workflow may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workflow" "example" {
  name = "Example workflow"
}

# Or by id:
data "oneuptime_workflow" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_email_secret_key` (String) The secret part of the Incoming Email trigger's address (workflow-<key>@<inbound email domain>). Anyone who has the address can start the workflow, so only people who can edit the workflow can read the key. Given to the workflow when its graph first has an Incoming Email trigger; set a new UUID to reset the address. Unique across all workflows.
- `is_archived` (Boolean) Archived workflows are hidden from the Workflows list and never run, from any trigger. Unarchiving restores them as they were.
- `is_enabled` (Boolean) Is this workflow enabled?
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.
- `webhook_secret_key` (String) The secret part of the Webhook trigger's URL (/workflow/trigger/<key>). Anyone who has the URL can start the workflow, so only people who can edit the workflow can read the key. Generated when the workflow is created; set a new value to reset the URL.

### Read-Only

- `archived_at` (String) When this workflow was archived. Empty while it is not archived.
- `created_at` (String) Date and Time when the object was created.
- `graph` (String) Workflow Graph in JSON. Ideally, create this via UI and not via API. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
