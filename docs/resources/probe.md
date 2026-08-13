---
page_title: "oneuptime_probe Resource - oneuptime"
subcategory: "Probes"
description: |-
  Manages custom probes. Deploy probes anywhere in the world and connect it to your project.
---

# oneuptime_probe (Resource)

Manages custom probes. Deploy probes anywhere in the world and connect it to your project.

## Example Usage

```terraform
resource "oneuptime_probe" "example" {
  key = "Example short text"
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  probe_version = jsonencode({
    "_type": "Version",
    "value": "1.0.0"
  })
  description = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
}
```

## Schema

### Required

- `key` (String) Probe key.
- `name` (String) Name object.
- `probe_version` (String) Version object.

### Optional

- `description` (String) Name object.
- `last_alive` (String) A date time object..
- `icon_file_id` (String) A unique identifier for an object, represented as a UUID..
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `should_auto_enable_probe_on_new_monitors` (Bool) Auto Enable Probe on New Monitors..
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `connection_status` (String) Connection Status of the Probe..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_probe.example <id>
```
