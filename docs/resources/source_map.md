---
page_title: "oneuptime_source_map Resource - oneuptime"
subcategory: "Other"
description: |-
  Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).
---

# oneuptime_source_map (Resource)

Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).

## Example Usage

```terraform
resource "oneuptime_source_map" "example" {
  service_id      = oneuptime_service.example.id
  service_version = "Example short text"
  bundle_path     = "This is an example of longer text content that might be stored in this field."
  content         = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `bundle_path` (String) Path or file name of the minified bundle this map was generated for (for example main.a8f1b2.js). Stack frames are matched against this by path suffix, so the file name alone is enough.
- `content` (String) The source map JSON (version 3) for this bundle.
- `service_id` (String) ID of the telemetry service this source map belongs to. The ID of a `oneuptime_service`.
- `service_version` (String) The release this source map belongs to. Must exactly match the service.version OpenTelemetry resource attribute sent with the telemetry.

### Optional

- `size_in_bytes` (Number) Size of the source map JSON in bytes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
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
