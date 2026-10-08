---
page_title: "oneuptime_source_map Data Source - oneuptime"
subcategory: "Other"
description: |-
  Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).
---

# oneuptime_source_map (Data Source)

Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one source map may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_source_map" "example" {
  service_id = oneuptime_service.example.id
}

# Or by id:
data "oneuptime_source_map" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `bundle_path` (String) Path or file name of the minified bundle this map was generated for (for example main.a8f1b2.js). Stack frames are matched against this by path suffix, so the file name alone is enough.
- `content` (String) The source map JSON (version 3) for this bundle.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `service_id` (String) ID of the telemetry service this source map belongs to. The ID of a `oneuptime_service`.
- `service_version` (String) The release this source map belongs to. Must exactly match the service.version OpenTelemetry resource attribute sent with the telemetry.
- `size_in_bytes` (Number) Size of the source map JSON in bytes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
