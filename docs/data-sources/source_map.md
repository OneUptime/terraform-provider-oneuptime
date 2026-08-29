---
page_title: "oneuptime_source_map Data Source - oneuptime"
subcategory: "Other"
description: |-
  Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute).
---

# oneuptime_source_map (Data Source)

Source maps uploaded for telemetry services. Used to resolve minified browser exception stack traces back to the original source code. Maps are matched to exceptions by service and release (the service.version OpenTelemetry resource attribute). Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_source_map" "by_name" {
  name = "example-source_map"
}

data "oneuptime_source_map" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_version` (String) The release this source map belongs to. Must exactly match the service.version OpenTelemetry resource attribute sent with the telemetry... Computed.
- `bundle_path` (String) Path or file name of the minified bundle this map was generated for (for example main.a8f1b2.js). Stack frames are matched against this by path suffix, so the file name alone is enough... Computed.
- `content` (String) The source map JSON (version 3) for this bundle.. Computed.
- `size_in_bytes` (Number) Size of the source map JSON in bytes.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
