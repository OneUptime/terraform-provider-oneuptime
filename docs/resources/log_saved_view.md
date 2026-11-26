---
page_title: "oneuptime_log_saved_view Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Save and reuse log explorer views, including the current filters, columns, sorting, and page size. Each view is listed on the page it was saved on: the Logs page, or one resource's Logs tab.
---

# oneuptime_log_saved_view (Resource)

Save and reuse log explorer views, including the current filters, columns, sorting, and page size. Each view is listed on the page it was saved on: the Logs page, or one resource's Logs tab.

## Example Usage

```terraform
resource "oneuptime_log_saved_view" "example" {
  name = "Example log saved view"
}
```

## Schema

### Required

- `name` (String) Friendly name for this saved log view.

### Optional

- `columns` (String) Selected log table columns for this saved view. A JSON value: write it with `jsonencode()`.
- `is_default` (Boolean) Whether this saved log view should be applied by default. Defaults to `false`.
- `page_size` (Number) Number of logs per page for this saved view. Defaults to `100`.
- `query` (String) Serialized log query for this saved view. A JSON value: write it with `jsonencode()`.
- `scope` (String) The page this saved log view is listed on: empty for the Logs page, or the resource page it was saved on, such as a Kubernetes cluster's Logs tab (kubernetes-cluster:<cluster id>).
- `sort_field` (String) Active sort field for this saved log view.
- `sort_order` (String) Sort order for this saved log view. Allowed values: `ASC`, `DESC`.
- `time_range` (String) Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this saved log view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this saved log view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log saved view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_saved_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_saved_view.example <id>
```
