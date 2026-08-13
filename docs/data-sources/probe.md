---
page_title: "oneuptime_probe Data Source - oneuptime"
subcategory: "Probes"
description: |-
  Manages custom probes. Deploy probes anywhere in the world and connect it to your project.
---

# oneuptime_probe (Data Source)

Manages custom probes. Deploy probes anywhere in the world and connect it to your project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_probe" "by_name" {
  name = "example-probe"
}

data "oneuptime_probe" "by_id" {
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
- `key` (String) Probe key. Computed.
- `description` (String) Name object. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `probe_version` (String) Version object. Computed.
- `last_alive` (String) A date time object.. Computed.
- `icon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `should_auto_enable_probe_on_new_monitors` (Bool) Auto Enable Probe on New Monitors.. Computed.
- `connection_status` (String) Connection Status of the Probe.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
