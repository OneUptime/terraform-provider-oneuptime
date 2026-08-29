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
  service_id = "123e4567-e89b-12d3-a456-426614174000"
  service_version = "Example short text"
  bundle_path = "This is an example of longer text content that might be stored in this field."
  content = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `service_id` (String) A unique identifier for an object, represented as a UUID..
- `service_version` (String) The release this source map belongs to. Must exactly match the service.version OpenTelemetry resource attribute sent with the telemetry...
- `bundle_path` (String) Path or file name of the minified bundle this map was generated for (for example main.a8f1b2.js). Stack frames are matched against this by path suffix, so the file name alone is enough...
- `content` (String) The source map JSON (version 3) for this bundle..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `size_in_bytes` (Number) Size of the source map JSON in bytes..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_source_map.example <id>
```
