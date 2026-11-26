---
page_title: "oneuptime_table_view Data Source - oneuptime"
subcategory: "Organization"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Data Source)

Table View is view settings for a table in a project. It contains columns, filters, and other settings.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one table view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_table_view" "example" {
  name = "Example table view"
}

# Or by id:
data "oneuptime_table_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `items_on_page` (Number) Items on page.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `table_id` (String) ID of the table this view is for.

### Read-Only

- `columns` (String) Which columns are shown, and in what order, for this table view. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `query` (String) Filters for this table view. A JSON value: read it with `jsondecode()`.
- `sort` (String) Sort for this table view. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
