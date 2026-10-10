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
  name        = "Example workflow"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `graph` (String) Workflow Graph in JSON. Ideally, create this via UI and not via API. A JSON value: write it with `jsonencode()`.
- `incoming_email_secret_key` (String) The secret part of the Incoming Email trigger's address (workflow-<key>@<inbound email domain>). Anyone who has the address can start the workflow, so only people who can edit the workflow can read the key. Given to the workflow when its graph first has an Incoming Email trigger; set a new UUID to reset the address. Unique across all workflows.
- `is_archived` (Boolean) Archived workflows are hidden from the Workflows list and never run, from any trigger. Unarchiving restores them as they were. Defaults to `false`.
- `is_enabled` (Boolean) Is this workflow enabled? Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `webhook_secret_key` (String) The secret part of the Webhook trigger's URL (/workflow/trigger/<key>). Anyone who has the URL can start the workflow, so only people who can edit the workflow can read the key. Generated when the workflow is created; set a new value to reset the URL.

### Read-Only

- `archived_at` (String) When this workflow was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing workflow by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_workflow.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_workflow.example <id>
```
