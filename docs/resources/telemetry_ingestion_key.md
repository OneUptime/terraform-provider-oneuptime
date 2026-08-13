---
page_title: "oneuptime_telemetry_ingestion_key Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Manage Telemetry Ingestion Keys for your project
---

# oneuptime_telemetry_ingestion_key (Resource)

Manage Telemetry Ingestion Keys for your project

## Example Usage

```terraform
resource "oneuptime_telemetry_ingestion_key" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `secret_key` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_telemetry_ingestion_key.example <id>
```
