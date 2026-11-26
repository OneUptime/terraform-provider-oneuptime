---
page_title: "oneuptime_source_map Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).
---

# oneuptime_source_map (Resource)

Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).

## Example Usage

```terraform
resource "oneuptime_source_map" "example" {
  service_id      = oneuptime_service.example.id
  service_version = "1.4.2"
  bundle_path     = "main.a8f1b2.js"
  content         = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `bundle_path` (String) Path or file name of the minified bundle this map was generated for (for example main.a8f1b2.js). Stack frames are matched against this by path suffix, so the file name alone is enough.
- `content` (String) The source map JSON (version 3) for this bundle.
- `service_id` (String) ID of the telemetry service this source map belongs to. The ID of a `oneuptime_service`.
- `service_version` (String) The release this source map belongs to. Must exactly match the service.version OpenTelemetry resource attribute sent with the telemetry.

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
- `size_in_bytes` (Number) Size of the source map JSON in bytes.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing source map by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_source_map.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_source_map.example <id>
```
