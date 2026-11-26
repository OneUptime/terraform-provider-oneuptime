---
page_title: "oneuptime_trace_saved_view Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Save and reuse traces explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_trace_saved_view (Resource)

Save and reuse traces explorer views, including the current search, filters, time range, and page size.

## Example Usage

```terraform
resource "oneuptime_trace_saved_view" "example" {
  name = "Example trace saved view"
}
```

## Schema

### Required

- `name` (String) Friendly name for this saved trace view.

### Optional

- `is_default` (Boolean) Whether this saved trace view should be applied by default. Defaults to `false`.
- `query` (String) Serialized traces explorer view state (search, filters, time range, page size) for this saved view. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this saved trace view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this saved trace view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing trace saved view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_trace_saved_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_trace_saved_view.example <id>
```
