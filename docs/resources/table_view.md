---
page_title: "oneuptime_table_view Resource - oneuptime"
subcategory: "Organization"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Resource)

Table View is view settings for a table in a project. It contains columns, filters, and other settings.

## Example Usage

```terraform
resource "oneuptime_table_view" "example" {
  name        = "Example table view"
  table_id    = "Monitor"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.
- `table_id` (String) ID of the table this view is for.

### Optional

- `columns` (String) Which columns are shown, and in what order, for this table view. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view. A JSON value: write it with `jsonencode()`.
- `items_on_page` (Number) Items on page. Defaults to `10`.
- `query` (String) Filters for this table view. A JSON value: write it with `jsonencode()`.
- `sort` (String) Sort for this table view. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing table view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_table_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_table_view.example <id>
```
